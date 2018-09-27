package config

import (
	"regexp"

	"github.com/mike-douglas/chaosproxy/middleware"
	"github.com/mike-douglas/chaosproxy/structs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type yamlRouteConfig struct {
	Pattern    string                 `yaml:"pattern"`
	Middleware []yamlMiddlewareConfig `yaml:"actions"`
}

func (r yamlRouteConfig) RouteConfig() (*structs.RouteConfig, error) {
	reg, err := regexp.Compile(r.Pattern)

	if err != nil {
		return nil, errors.Wrap(err, "Could not compile pattern regex")
	}

	var mw = make([]structs.MiddlewareConfig, len(r.Middleware))

	for i, m := range r.Middleware {
		mw[i] = *m.MiddlewareConfig()
	}

	var rc = structs.RouteConfig{
		CompiledPattern: reg,
		Middleware:      mw,
		Handler:         middleware.CreateHandler(mw),
	}

	return &rc, nil
}

type yamlMiddlewareConfig struct {
	Action string            `yaml:"name"`
	Params map[string]string `yaml:"params"`
}

func (m yamlMiddlewareConfig) MiddlewareConfig() *structs.MiddlewareConfig {
	return &structs.MiddlewareConfig{
		Action: m.Action,
		Params: m.Params,
	}
}

// BuildFromYaml ...
func BuildFromYaml(source []byte, c *structs.Config) error {
	var yamlRoutes = []yamlRouteConfig{}
	var routes = []structs.RouteConfig{}
	var route *structs.RouteConfig
	var err error

	err = yaml.Unmarshal(source, &yamlRoutes)

	if err != nil {
		return errors.Wrap(err, "Could not parse yaml file")
	}

	for i, yr := range yamlRoutes {
		route, err = yr.RouteConfig()

		if err != nil {
			return errors.Wrapf(err, "Could not build config for route %d", i)
		}

		routes = append(routes, *route)
	}

	c.Routes = routes

	return nil
}
