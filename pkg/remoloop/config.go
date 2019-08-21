package remoloop

import "github.com/kazukousen/remoloop/pkg/remoloop/client"

// Config ...
type Config struct {
	Client client.Config `yaml:"client"`
}
