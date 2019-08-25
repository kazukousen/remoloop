package remoloop

import (
	"github.com/go-kit/kit/log"
	"github.com/kazukousen/remoloop/pkg/remoloop/client"
	"github.com/kazukousen/remoloop/pkg/remoloop/server"
)

func New(cfg *Config, logger log.Logger) (server.Server, error) {
	client, err := client.New(logger, cfg.Client)
	if err != nil {
		return nil, err
	}
	return server.New(logger, cfg.Server, client)
}
