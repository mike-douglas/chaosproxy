package config

import (
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
	route, err := buildRouteConfig()

	if err != nil {
		return nil, errors.Wrap(err, "Could not build config")
	}

	return &structs.Config{
		Verbose: false,
		Routes:  []structs.RouteConfig{*route},
	}, nil
}

func buildRouteConfig() (*structs.RouteConfig, error) {
	route := structs.RouteConfig{
		Pattern: "abmcscholar.gov",
		Middleware: []structs.MiddlewareConfig{
			{Action: "sleep", Params: map[string]string{"seconds": "1"}},
			{Action: "randsleep", Params: map[string]string{"from": "1", "to": "4"}},
			{Action: "sleep", Params: map[string]string{"seconds": "1"}},
		},
	}

	if r, err := regexp.Compile(route.Pattern); err == nil {
		route.CompiledPattern = r
	} else {
		return nil, errors.Wrap(err, "Could not compile route pattern")
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

	return &route, nil
}
