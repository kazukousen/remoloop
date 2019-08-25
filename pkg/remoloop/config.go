package remoloop

import (
	"github.com/kazukousen/remoloop/pkg/remoloop/client"
	"github.com/kazukousen/remoloop/pkg/remoloop/server"
)

// Config ...
type Config struct {
	Client client.Config `yaml:"client"`
	Server server.Config `yaml:"server"`
}
