package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
			devices := []device{}
			if err := c.request(context.Background(), resourceDevices, &devices); err != nil {
				level.Error(c.logger).Log("msg", "could not request", "path", resourceDevices, "error", err)
				continue
			}
			for _, device := range devices {
				level.Debug(c.logger).Log("msg", "get devices", "id", device.ID, "name", device.Name)
			}
		case <-c.exit:
			return
		}
	}
}

func (c client) request(ctx context.Context, resource resource, dst interface{}) error {
	req, err := http.NewRequest(http.MethodGet, c.host+string(resource), nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	ch := make(chan error, 1)
	go func() {
		ch <- c.doRequest(req, dst)
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

func (c client) doRequest(req *http.Request, dst interface{}) (err error) {
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// level.Debug(c.logger).Log("msg", "json data", "data", string(b))
	err = json.Unmarshal(b, dst)
	return
}
