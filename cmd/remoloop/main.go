package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kazukousen/remoloop/pkg/helpers"
	"github.com/kazukousen/remoloop/pkg/remoloop"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	if len(os.Args) == 0 {
		level.Error(logger).Log("msg", "require argument")
	}
	filename := os.Args[1]
	cfg := &remoloop.Config{}
	err := helpers.LoadYAML(filename, cfg)
	if err != nil {
		level.Error(logger).Log("msg", "could not load to yaml", "error", err)
		return
	}

	level.Info(logger).Log("msg", "start processing", "config file", filename)

	server, err := remoloop.New(cfg, logger)
	if err != nil {
		level.Error(logger).Log("msg", "could not initialize server", "error", err)
		return
	}

	// wait signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	level.Info(logger).Log("msg", "signal recieved, then shutting down", "signal", <-quit)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	server.Stop(ctx)
}
