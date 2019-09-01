package remoloop

import (
	"github.com/go-kit/kit/log"
	"github.com/kazukousen/remoloop/pkg/remoloop/client"
	"github.com/kazukousen/remoloop/pkg/remoloop/server"
)

type RemoLoop interface {
	Stop()
}

type remoloop struct {
	client client.Client
	server server.Server
}

func New(cfg *Config, logger log.Logger) (RemoLoop, error) {
	client, err := client.New(logger, cfg.Client)
	if err != nil {
		return nil, err
	}

	server, err := server.New(logger, cfg.Server, client)
	if err != nil {
		return nil, err
	}

	return &remoloop{
		client: client,
		server: server,
	}, nil
}

func (r remoloop) Stop() {
	r.server.Stop()
	r.client.Stop()
}
