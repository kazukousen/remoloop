package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kazukousen/remoloop/pkg/helpers"
	"github.com/kazukousen/remoloop/pkg/remoloop"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	stop := signalHandler(logger)

	if len(os.Args) < 2 {
		level.Error(logger).Log("msg", "require argument")
		return
	}
	filename, cfg := os.Args[1], &remoloop.Config{}
	if err := helpers.LoadYAML(filename, cfg); err != nil {
		level.Error(logger).Log("msg", "could not load to yaml", "error", err)
		return
	}

	server, err := remoloop.New(cfg, logger)
	if err != nil {
		level.Error(logger).Log("msg", "could not initialize server", "error", err)
		return
	}

	level.Info(logger).Log("msg", "start processing", "config file", filename)
	<-stop
	server.Stop()
	level.Info(logger).Log("msg", "see you.")
}

func signalHandler(logger log.Logger) <-chan struct{} {
	stop := make(chan struct{}, 1)
	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
		level.Info(logger).Log("msg", "signal received, then shutting down", "signal", <-quit)
		close(stop)
		level.Info(logger).Log("msg", "second signal received. exit directly", "signal", <-quit)
		os.Exit(1)
	}()
	return stop
}
