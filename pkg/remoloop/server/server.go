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
	done   chan struct{}
}

// New returns Server and http Server listen and serve.
func New(logger log.Logger, cfg Config, client client.Client) (Server, error) {
	path, port := "", 4100
	if cfg.RootPath != "" {
		path = cfg.RootPath
		if path[0] != '/' {
			path = "/" + path
		}
	}
	if cfg.Port != 0 {
		port = cfg.Port
	}

	s := &server{
		logger: log.With(logger, "component", "server"),
		done:   make(chan struct{}, 1),
		client: client,
	}
	s.server = &http.Server{
		ReadHeaderTimeout: 1 * time.Minute,
		Addr:              ":" + strconv.Itoa(port),
		Handler:           s,
	}

	if path != "" {
		s.server.Handler = s.rewriteRootPath(path, s)
	} else {
		s.server.Handler = s
	}

	go s.serve()

	return s, nil
}

func (s server) serve() {
	s.server.ListenAndServe()
	close(s.done)
}

func (s server) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		level.Error(s.logger).Log("msg", "failed to shutdown server", "error", err)
		return
	}
	<-s.done
	level.Info(s.logger).Log("msg", "success to stop server")
	s.client.Stop(ctx)
}
