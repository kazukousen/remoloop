package client_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi"

	"github.com/go-kit/kit/log"
	"github.com/kazukousen/remoloop/pkg/helpers"
	"github.com/kazukousen/remoloop/pkg/remoloop/api"
	"github.com/kazukousen/remoloop/pkg/remoloop/client"
)

type controller struct {
	mu    *sync.Mutex
	inc   int
	token string
}

func newController() *controller {
	return &controller{
		mu: &sync.Mutex{},
	}
}

func (c *controller) Inc(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.inc++
	w.WriteHeader(http.StatusOK)
}

func (c *controller) Me(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if bt := r.Header.Get("Authorization"); strings.HasPrefix(bt, "Bearer ") {
		token := bt[len("Bearer "):]
		c.token = token
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`
{
	"nickname": "Alice"
}
		`))
}

func TestClient_Get(t *testing.T) {
	r, c := chi.NewRouter(), newController()
	r.Get("/1/users/me", c.Me)
	r.Get("/foo", c.Inc)
	srv := httptest.NewServer(r)
	defer srv.Close()

	cfg := client.Config{
		Host: srv.URL,
		HTTPConfig: helpers.HTTPClientConfig{
			BearerToken: "fake-token",
		},
	}
	client, err := client.New(log.NewLogfmtLogger(ioutil.Discard), cfg)
	if err != nil {
		t.Errorf("could not initialize client: %+v", err)
		return
	}
	if c.token != "fake-token" {
		t.Errorf("not equal got %s, but want %s", c.token, "fake-token")
		return
	}

	for i := 0; i < 10; i++ {
		client.Get(context.Background(), api.Resource("/foo"), ioutil.Discard)
	}
	if c.inc != 10 {
		t.Errorf("not equal got %d, but want %d", c.inc, 10)
		return
	}
}
