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

// RouteConfig ...
type RouteConfig struct {
	Pattern         string             `yaml:"pattern"`
	Middleware      []MiddlewareConfig `yaml:"actions"`
	CompiledPattern *regexp.Regexp
	Handler         http.Handler
}

func (rc *RouteConfig) String() string {
	return fmt.Sprintf("Route (Handler: %v; Pattern: %v)", rc.Handler, rc.CompiledPattern)
}

// MiddlewareConfig ...
type MiddlewareConfig struct {
	Action string            `yaml:"name"`
	Params map[string]string `yaml:"params"`
}
