package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kazukousen/remoloop/pkg/helpers"
	"github.com/kazukousen/remoloop/pkg/remoloop/api"
	"golang.org/x/xerrors"
)

// Client represents API Client.
type Client interface {
	Stop(ctx context.Context)
	Get(ctx context.Context, resource api.Resource, w io.Writer)
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
	if cfg.Host != "" {
		host = cfg.Host
	}
	c := &client{
		host:   host,
		logger: log.With(logger, "component", "client", "host", host),
		exit:   make(chan struct{}, 1),
		done:   make(chan struct{}, 1),
		client: helpers.NewHTTPClient(cfg.HTTPConfig),
	}
	c.client.Transport = newRateLimitRoundTripper(c.client.Transport)

	// ping
	buf := &bytes.Buffer{}
	if err := c.request(context.Background(), api.ResourceUsersMe, buf); err != nil {
		return nil, xerrors.Errorf("could not request: %w", err)
	}
	dst := &api.Me{}
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

type rateLimitRoundTripper struct {
	rt        http.RoundTripper
	mu        *sync.Mutex
	reset     time.Time
	remaining int
}

func newRateLimitRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &rateLimitRoundTripper{
		rt: rt,
		mu: &sync.Mutex{},
	}
}

func (lrt *rateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !lrt.reset.IsZero() && lrt.reset.Sub(time.Now()) > 0 && lrt.remaining < 1 {
		return nil, xerrors.Errorf("rate limit quota. the time until the next reset: %s", lrt.reset)
	}

	// RoundTrip
	res, err := lrt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if err := lrt.setRateLimit(res); err != nil {
		return nil, err
	}

	return res, err
}

func (lrt rateLimitRoundTripper) setRateLimit(res *http.Response) error {
	lrt.mu.Lock()
	defer lrt.mu.Unlock()
	// Number of remaining requests
	remaining, _ := strconv.Atoi(res.Header.Get("X-Rate-Limit-Remaining"))
	lrt.remaining = remaining

	// Time until the next reset
	reset, _ := strconv.ParseInt(res.Header.Get("X-Rate-Limit-Reset"), 10, 64)
	lrt.reset = time.Unix(reset, 0)
	return nil
}
