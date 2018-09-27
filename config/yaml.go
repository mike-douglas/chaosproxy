package config

import (
	"fmt"

	"github.com/mike-douglas/chaosproxy/structs"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type routeYaml struct {
	pattern    string           `yaml:"string"`
	middleware []middlewareYaml `yaml:"actions,flow"`
}

type middlewareYaml struct {
	action string            `yaml:"action"`
	params map[string]string `yaml:"params"`
}

// BuildFromYaml ...
func BuildFromYaml(source []byte, c *structs.Config) error {
	var yamlRoutes = routeYaml{}

	err := yaml.Unmarshal(source, &yamlRoutes)

	if err != nil {
		return errors.Wrap(err, "Could not parse yaml file")
	}

	// Use yamlroutes to populate c.Routes

	for i, yr := range yamlRoutes {
		fmt.Printf("%v\n", route)
		err = buildRouteConfig(&route)

		fmt.Printf("%v\n", route)

		if err != nil {
			return errors.Wrapf(err, "Could not build config for route %d", i)
		}
	}

	return nil
}
