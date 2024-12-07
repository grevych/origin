// Description: This file is the focal point of configuration that needs passed
// to various parts of the service.
package config

import (
	"context"

	"github.com/grevych/gobox/pkg/cfg"
	"github.com/grevych/gobox/pkg/log"
)

// Config tracks config needed for origin
type Config struct {
	listenHost      string `yaml:"listenHost"`
	privateHTTPPort int    `yaml:"privateHTTPPort"`
	publicHTTPPort  int    `yaml:"publicHTTPPort"`
	gRPCPort        int    `yaml:"gRPCPort"`
}

// Load returns a new Config type that has been loaded in accordance to the environment
// that the service was deployed in, with all necessary tweaks made before returning.
func Load(ctx context.Context) (*Config, error) {
	// NOTE: Defaults should generally be set in the config
	// override jsonnet file: deployments/{{ .Config.Name }}/{{ .Config.Name }}.config.jsonnet
	c := Config{
		// Defaults to [::]/0.0.0.0 which will broadcast to all reachable
		// IPs on a server on the given port for the respective service.
		listenHost:      "",
		privateHTTPPort: 8000,
		publicHTTPPort:  8080,
		gRPCPort:        5000,
	}

	// Attempt to load a local config file on top of the defaults
	if err := cfg.Load("config.yaml", &c); err != nil {
		return nil, err
	}

	log.Info(ctx, "Configuration data of the application:\n", &c)

	return &c, nil
}

// MarshalLog can be used to write config to log
func (c *Config) MarshalLog(addfield func(key string, value interface{})) {
}

func (c *Config) ListenHost() string {
	return c.listenHost
}

func (c *Config) PrivateHTTPPort() int {
	return c.privateHTTPPort
}

func (c *Config) PublicHTTPPort() int {
	return c.publicHTTPPort
}

func (c *Config) GRPCPort() int {
	return c.gRPCPort
}
