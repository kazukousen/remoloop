package client

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kazukousen/remoloop/pkg/helpers"
)

// Client ...
type Client interface {
	Stop()
}

type client struct {
	logger log.Logger
	cfg    Config
	host   string
	client *http.Client
	exit   chan struct{}
	done   chan struct{}
}

// NewClient ...
func NewClient(logger log.Logger, cfg Config) (Client, error) {
	host := "https://api.nature.global"
	c := &client{
		host:   host,
		logger: log.With(logger, "component", "client", "host", host),
		exit:   make(chan struct{}, 1),
		done:   make(chan struct{}, 1),
	}
	c.client = helpers.NewHTTPClient(cfg.HTTPClientConfig)
	dst := &usersMe{}
	if err := c.request(context.Background(), resourceUsersMe, dst); err != nil {
		return nil, err
	}
	level.Debug(c.logger).Log("msg", "sucess to authorization", "nickname", dst.Nickname)

	go c.work()

	return c, nil
}

func (c client) Stop() {
	close(c.exit)
	<-c.done
}
