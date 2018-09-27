package config

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/mike-douglas/chaosproxy/middleware"
	"github.com/mike-douglas/chaosproxy/structs"
	"github.com/pkg/errors"
)

// BuildConfig ...
func BuildConfig() (*structs.Config, error) {
	var route = structs.RouteConfig{
		Pattern: "abmcscholar.gov",
		Middleware: []structs.MiddlewareConfig{
			{Action: "sleep", Params: map[string]string{"seconds": "1"}},
			{Action: "randsleep", Params: map[string]string{"from": "1", "to": "4"}},
			{Action: "sleep", Params: map[string]string{"seconds": "1"}},
		},
	}

	err := buildRouteConfig(&route)

	if err != nil {
		return nil, errors.Wrap(err, "Could not build config")
	}

	return &structs.Config{
		Verbose: false,
		Routes:  []structs.RouteConfig{route},
	}, nil
}

func buildRouteConfig(route *structs.RouteConfig) error {
	if r, ok := regexp.Compile(route.Pattern); ok == nil {
		route.CompiledPattern = r.Copy()
	} else {
		return errors.Wrap(ok, "Could not compile route pattern")
	}

	var handler http.Handler
	firstMiddleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if url, ok := url.Parse(r.RequestURI); ok == nil {
			httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})

	for i, c := range route.Middleware {
		mw := middleware.GetMiddlewareFunction(c.Action)
		if i == 0 {
			handler = mw(c, firstMiddleware)
		} else {
			handler = mw(c, handler)
		}
	}

	route.Handler = handler

	fmt.Printf("%v\n", route)

	return nil
}
