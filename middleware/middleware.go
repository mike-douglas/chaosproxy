package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
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
	case "noop":
		return noopMiddleware
	default:
		return noopMiddleware
	}
}

func noopMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Noop")
		next.ServeHTTP(w, r)
	})
}

func sleepMiddleware(c structs.MiddlewareConfig, next http.Handler) http.Handler {
	seconds, err := strconv.Atoi(c.Params["seconds"])

	if err != nil {
		// Dieeee
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Sleeeeeep...")
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

	randomSleep := time.Duration(rand.Intn(to)+from) * 1000 * time.Millisecond
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Random Sleeeeeep %v...\n", randomSleep)
		time.Sleep(randomSleep)
		next.ServeHTTP(w, r)
	})
}
