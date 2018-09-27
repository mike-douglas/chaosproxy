package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// MiddlewareConfig ...
type MiddlewareConfig struct {
	action string
	params map[string]string
}

func noopMiddleware(c MiddlewareConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Noop")
		next.ServeHTTP(w, r)
	})
}

func sleepMiddleware(c MiddlewareConfig, next http.Handler) http.Handler {
	seconds, err := strconv.Atoi(c.params["seconds"])

	if err != nil {
		// Dieeee
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Sleeeeeep...")
		time.Sleep(time.Duration(seconds) * 1000 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}

func randomSleepMiddleware(c MiddlewareConfig, next http.Handler) http.Handler {
	from, err := strconv.Atoi(c.params["from"])

	if err != nil {
		// Dieeee
	}

	to, err := strconv.Atoi(c.params["to"])

	if err != nil {
		// Dieeee
	}

	randomSleep := time.Duration(rand.Intn(to)+from) * 1000 * time.Millisecond
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Random Sleeeeeep %v...\n", randomSleep)
		time.Sleep(randomSleep)
		next.ServeHTTP(w, r)
	})
}

func getMiddlewareFunction(action string) func(MiddlewareConfig, http.Handler) http.Handler {
	switch action {
	case "sleep":
		return sleepMiddleware
	case "randsleep":
		return randomSleepMiddleware
	default:
		return noopMiddleware
	}
}

func proxyRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request: %s\n", r.RequestURI)

	config := []MiddlewareConfig{
		{action: "sleep", params: map[string]string{"seconds": "1"}},
		{action: "randsleep", params: map[string]string{"from": "1", "to": "4"}},
		{action: "sleep", params: map[string]string{"seconds": "1"}},
	}

	if url, ok := url.Parse(r.RequestURI); ok == nil {
		proxy := httputil.NewSingleHostReverseProxy(url)
		var middleware http.Handler

		for i, c := range config {
			mw := getMiddlewareFunction(c.action)
			if i == 0 {
				middleware = mw(c, proxy)
			} else {
				middleware = mw(c, middleware)
			}
		}

		middleware.ServeHTTP(w, r)
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", proxyRequestHandler)

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
