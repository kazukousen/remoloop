package client

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/kazukousen/remoloop/pkg/remoloop/api"

	"github.com/go-kit/kit/log/level"
	"golang.org/x/xerrors"
)

func (c client) Get(ctx context.Context, resource api.Resource, w io.Writer) {
	if err := c.request(ctx, resource, w); err != nil {
		level.Error(c.logger).Log("msg", "object", "error", err)
	}
}

func (c client) request(ctx context.Context, resource api.Resource, w io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, c.host+resource.String(), nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	ch := make(chan error, 1)
	go func() {
		ch <- c.doRequest(req, w)
	}()

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		<-ch
		err = ctx.Err()
		return err
	}

	return nil
}

func (c client) doRequest(req *http.Request, w io.Writer) (err error) {
	res, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			err = e
		}
	}()
	if res.StatusCode != http.StatusOK {
		err = xerrors.Errorf("http status code: %s", res.StatusCode)
		return
	}
	_, err = io.Copy(w, res.Body)
	return
}
