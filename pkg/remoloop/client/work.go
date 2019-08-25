package client

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/kazukousen/remoloop/pkg/remoloop/api"

	"github.com/go-kit/kit/log/level"
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
			if err := c.request(context.Background(), api.ResourceAppliance, buf); err != nil {
				level.Error(c.logger).Log("msg", "could not request", "resource", api.ResourceAppliance, "error", err)
				continue
			}
			appliances := []api.Appliance{}
			if err := json.Unmarshal(buf.Bytes(), &appliances); err != nil {
				level.Error(c.logger).Log("msg", "failed to decode json", "resource", api.ResourceAppliance, "error", err)
				continue
			}
		case <-c.exit:
			return
		}
	}
}
