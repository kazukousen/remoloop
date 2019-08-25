package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/kazukousen/remoloop/pkg/remoloop/client"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Server represents http.Server.
type Server interface {
	Stop(ctx context.Context)
}

type server struct {
	logger log.Logger
	server *http.Server
	client client.Client
	done   chan error
}

// New returns Server and http Server listen and serve.
func New(logger log.Logger, cfg Config, client client.Client) (Server, error) {
	path, port := "/", 4100
	if cfg.RootPath != "" {
		path = cfg.RootPath
		if path[0] != '/' {
			path = "/" + path
		}
	}
	if cfg.Port != 0 {
		port = cfg.Port
	}

	handler := &handler{path: path}

	s := &server{
		logger: log.With(logger, "component", "server"),
		done:   make(chan error, 1),
		client: client,
		server: &http.Server{
			ReadHeaderTimeout: 1 * time.Minute,
			Addr:              ":" + strconv.Itoa(port),
			Handler:           handler,
		},
	}

	go s.run()

	return s, nil
}

func (s server) run() {
	s.done <- s.server.ListenAndServe()
}

func (s server) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		level.Error(s.logger).Log("msg", "failed to shutdown server", "error", err)
		return
	}
	level.Info(s.logger).Log("msg", "success to stop server")
	s.client.Stop(ctx)
}
