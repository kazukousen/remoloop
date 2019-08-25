package client

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"io"
)

func (c client) Get(resource resource, w io.Writer) {
	if err := c.request(context.Background(), resource, w); err != nil {
		level.Error(c.logger).Log("msg", "object", "error", err)
	}
}
