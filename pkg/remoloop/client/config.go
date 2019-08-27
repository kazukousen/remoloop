package client

import "github.com/kazukousen/remoloop/pkg/helpers"

// Config ...
type Config struct {
	Host       string                   `yaml:"host,omitempty"`
	HTTPConfig helpers.HTTPClientConfig `yaml:",inline"`
}
