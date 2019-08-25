package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"golang.org/x/xerrors"

	"github.com/go-kit/kit/log"
	"github.com/kazukousen/remoloop/pkg/helpers"
)

// Client represents API Client.
type Client interface {
	Stop(ctx context.Context)
	Get(resource resource, w io.Writer)
}

type client struct {
	logger log.Logger
	cfg    Config
	host   string
	client *http.Client
	exit   chan struct{}
	done   chan struct{}
}

// New returns Client and stand up worker goroutine.
func New(logger log.Logger, cfg Config) (Client, error) {
	host := "https://api.nature.global"
	c := &client{
		host:   host,
		logger: log.With(logger, "component", "client"),
		exit:   make(chan struct{}, 1),
		done:   make(chan struct{}, 1),
	}
	c.client = helpers.NewHTTPClient(cfg.HTTPClientConfig)

	// ping
	buf := &bytes.Buffer{}
	if err := c.request(context.Background(), resourceUsersMe, buf); err != nil {
		return nil, err
	}
	dst := &usersMe{}
	if err := json.Unmarshal(buf.Bytes(), dst); err != nil {
		return nil, xerrors.New("could not connect API")
	}
	level.Info(c.logger).Log("msg", "sucess to authorization", "nickname", dst.Nickname)

	go c.work()

	return c, nil
}

func (c client) Stop(ctx context.Context) {
	close(c.exit)
	select {
	case <-c.done:
		level.Info(c.logger).Log("msg", "success to stop client")
	case <-ctx.Done():
	}
}
