package structs

import (
	"net/http"
	"regexp"
)

// Config ...
type Config struct {
	Verbose bool
	Routes  []RouteConfig
}

// RouteConfig ...
type RouteConfig struct {
	Pattern         string
	Middleware      []MiddlewareConfig
	CompiledPattern *regexp.Regexp
	Handler         http.Handler
}

// MiddlewareConfig ...
type MiddlewareConfig struct {
	Action string
	Params map[string]string
}
