package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/mike-douglas/chaosproxy/structs"
)

// GetMiddlewareFunction ...
func GetMiddlewareFunction(action string) func(structs.MiddlewareConfig, http.Handler) http.Handler {
	switch action {
	case "sleep":
		return sleepMiddleware
	case "randsleep":
		return randomSleepMiddleware
	case "httperror":
		return errorMiddleware
	case "noop":
		return noopMiddleware
	default:
		return noopMiddleware
	}
}

// CreateHandler creates a single handler for handling proxied requests
func CreateHandler(middleware []structs.MiddlewareConfig) http.Handler {
	var handler http.Handler
	lastMiddleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if url, ok := url.Parse(r.RequestURI); ok == nil {
			httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})

	for i, c := range middleware {
		mw := GetMiddlewareFunction(c.Action)
		if i == 0 {
			handler = mw(c, lastMiddleware)
		} else {
			handler = mw(c, handler)
		}
	}

	return handler
}

func noopMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(c)
		next.ServeHTTP(w, r)
	})
}

func errorMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	httpError, err := strconv.Atoi(c.Params["code"])

	if err != nil {
		// Error
	}

	httpMessage := c.Params["status"]

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(c)
		http.Error(w, httpMessage, httpError)
	})
}

func sleepMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	seconds, err := strconv.Atoi(c.Params["seconds"])

	if err != nil {
		// Dieeee
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(c)
		time.Sleep(time.Duration(seconds) * 1000 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}

func randomSleepMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	from, err := strconv.Atoi(c.Params["from"])

	if err != nil {
		// Dieeee
	}

	to, err := strconv.Atoi(c.Params["to"])

	if err != nil {
		// Dieeee
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(c)
		time.Sleep(time.Duration(rand.Intn(to)+from) * 1000 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}
