package client

import "github.com/kazukousen/remoloop/pkg/helpers"

// Config ...
type Config struct {
	helpers.HTTPClientConfig `yaml:",inline"`
}
