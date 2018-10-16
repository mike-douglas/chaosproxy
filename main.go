package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/kr/mitm"

	"github.com/gorilla/handlers"
	"github.com/mike-douglas/chaosproxy/config"
	"github.com/mike-douglas/chaosproxy/structs"
)

func buildHandlerFromConfig(config *structs.Config, upstream http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Verbose {
			fmt.Printf("Request: %v\n", r)
		}

		if url, err := url.Parse(r.RequestURI); err == nil {
			for _, route := range config.Routes {
				fmt.Printf("%v\n", route)
				if route.CompiledPattern.MatchString(url.String()) {
					route.Handler.ServeHTTP(w, r)
					return
				}
			}
			upstream.ServeHTTP(w, r)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	}
}

func main() {
	f := flag.String("config", "", "Configuration YAML file for proxy")

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

	p := &mitm.Proxy{
		CA: nil,
		Wrap: func(upstream http.Handler) http.Handler {
			return buildHandlerFromConfig(&c, upstream)
		},
	}

	if c.Verbose {
		fmt.Println(c)
	}

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, p))
}
