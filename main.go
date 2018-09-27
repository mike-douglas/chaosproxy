package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mike-douglas/chaosproxy/config"
	"github.com/mike-douglas/chaosproxy/structs"
)

func buildHandlerFromConfig(config *structs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request: %s\n", r.RequestURI)

		if url, err := url.Parse(r.RequestURI); err == nil {
			for _, route := range config.Routes {
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
	c, err := config.BuildConfig()

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", buildHandlerFromConfig(c))

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
