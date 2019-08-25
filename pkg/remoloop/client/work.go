package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/log/level"
	"golang.org/x/xerrors"
)

func (c client) work() {
	t := time.NewTicker(10 * time.Second)
	defer func() {
		t.Stop()
		close(c.done)
	}()

	for {
		select {
		case <-t.C:
			buf := &bytes.Buffer{}
			if err := c.request(context.Background(), resourceDevices, buf); err != nil {
				level.Error(c.logger).Log("msg", "could not request", "resource", resourceDevices, "error", err)
				continue
			}
			devices := []device{}
			if err := json.Unmarshal(buf.Bytes(), &devices); err != nil {
				level.Error(c.logger).Log("msg", "failed decode", "resource", resourceDevices, "error", err)
				continue
			}
			/*
				for _, device := range devices {
					level.Debug(c.logger).Log("msg", "get devices", "id", device.ID, "name", device.Name)
				}
			*/
		case <-c.exit:
			return
		}
	}
}

func (c client) request(ctx context.Context, resource resource, w io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, c.host+string(resource), nil)
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
