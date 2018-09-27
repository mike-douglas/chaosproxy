package structs

import (
	"fmt"
	"net/http"
	"regexp"
)

// Config ...
type Config struct {
	Verbose bool
	Routes  []RouteConfig
}

func (c Config) String() string {
	return fmt.Sprintf("Config(Verbose=%v; Routes=%v", c.Verbose, c.Routes)
}

// RouteConfig ...
type RouteConfig struct {
	Pattern         string
	Middleware      []MiddlewareConfig
	CompiledPattern *regexp.Regexp
	Handler         http.Handler
}

func (rc RouteConfig) String() string {
	return fmt.Sprintf("Route(Pattern=%v; Middleware=%v)", rc.CompiledPattern, rc.Middleware)
}

// MiddlewareConfig ...
type MiddlewareConfig struct {
	Action string
	Params map[string]string
}

func (mw MiddlewareConfig) String() string {
	return fmt.Sprintf("Middleware(Action=%v; Params=%v)", mw.Action, mw.Params)
}
