package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func proxyRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request: %s\n", r.RequestURI)

	if url, ok := url.Parse(r.RequestURI); ok == nil {
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", proxyRequestHandler)

	//  Read through config file, add patterns to the router?

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
