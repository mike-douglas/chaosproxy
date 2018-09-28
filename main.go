package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mike-douglas/chaosproxy/config"
	"github.com/mike-douglas/chaosproxy/structs"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest, err := net.DialTimeout("tcp", r.Host, 10*time.Second)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)

	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	client, _, err := hijacker.Hijack()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(dest, client)
	go transfer(client, dest)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func buildHandlerFromConfig(config *structs.Config) http.HandlerFunc {
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
			httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
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

	r := mux.NewRouter()
	h := buildHandlerFromConfig(&c)

	r.HandleFunc("/", h).Methods("POST", "GET", "PUT", "PATCH", "DELETE", "CONNECT")

	if c.Verbose {
		fmt.Println(c)
	}

	fmt.Printf("Listening on %s\n", fmt.Sprintf(":%v", *p))
	fmt.Println(*pem)

	if len(*pem) > 0 && len(*key) > 0 {
		fmt.Println("Starting SSL Server")
		server := &http.Server{
			Addr: fmt.Sprintf(":%v", *p),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodConnect {
					handleTunneling(w, r)
				} else {
					h(w, r)
				}
			}),
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		server.ListenAndServeTLS(*pem, *key)
	} else {
		http.ListenAndServe(fmt.Sprintf(":%v", *p), handlers.LoggingHandler(os.Stdout, r))
	}
}
