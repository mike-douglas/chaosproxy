package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kr/mitm"

	"github.com/mike-douglas/chaosproxy/config"
	"github.com/mike-douglas/chaosproxy/structs"
)

func buildHandlerFromConfig(config *structs.Config, upstream http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Verbose {
			fmt.Printf("Request: %v\n", r)
		}

		if url, err := url.Parse(r.RequestURI); err == nil {
			if strings.ToLower(url.Scheme) == "https" {
				// Skip HTTPs processing for now...
				upstream.ServeHTTP(w, r)
			} else {
				for _, route := range config.Routes {
					fmt.Printf("%v\n", route)
					if route.CompiledPattern.MatchString(url.String()) {
						route.Handler.ServeHTTP(w, r)
						return
					}
				}
				upstream.ServeHTTP(w, r)
			}
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	}
}

func main() {
	f := flag.String("config", "", "Configuration YAML file for proxy")
	p := flag.Int("port", 8080, "The port to listen on")
	pem := flag.String("pem", "", "SSL Cert")
	key := flag.String("key", "", "SSL Key")

	flag.Parse()

	d, err := ioutil.ReadFile(*f)

	if err != nil {
		panic(err)
	}

	var c = structs.Config{
		Verbose: true,
	}

	err = config.BuildFromYaml(d, &c)

	if err != nil {
		panic(err)
	}

	cert, err := tls.LoadX509KeyPair(*pem, *key)

	if err != nil {
		panic(err)
	}

	proxy := &mitm.Proxy{
		CA: &cert,
		TLSClientConfig: &tls.Config{
			KeyLogWriter:       os.Stdout,
			InsecureSkipVerify: true,
		},
		Wrap: func(upstream http.Handler) http.Handler {
			return buildHandlerFromConfig(&c, upstream)
		},
	}

	if c.Verbose {
		fmt.Println(c)
	}

	fmt.Printf("Port = %d\nSSL cert = %s\nSSL key = %s\n", *p, *pem, *key)
	http.ListenAndServe(fmt.Sprintf(":%d", *p), proxy)
}
