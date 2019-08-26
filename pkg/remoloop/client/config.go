package client

import "github.com/kazukousen/remoloop/pkg/helpers"

// Config ...
type Config struct {
	host                     string `yaml:"host"`
	helpers.HTTPClientConfig `yaml:",inline"`
}
