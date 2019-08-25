package helpers

// TODO:

import (
	"net/http"
	"time"
)

// HTTPClientConfig represents HTTP Client configuration.
type HTTPClientConfig struct {
	BearerToken string `yaml:"bearer_token"`
}

// NewHTTPClient returns http.Client.
func NewHTTPClient(cfg HTTPClientConfig) *http.Client {
	rt := &http.Transport{
		IdleConnTimeout:     5 * time.Minute,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{
		Transport: rt,
	}
	client.Transport = newBearerAuthRoudTripper(rt, cfg.BearerToken)
	return client
}

// bearerAuthRoundTripper implements http.RoundTripper.
type bearerAuthRoundTripper struct {
	bearerToken string
	rt          http.RoundTripper
}

func newBearerAuthRoudTripper(rt http.RoundTripper, token string) http.RoundTripper {
	return &bearerAuthRoundTripper{
		bearerToken: token,
		rt:          rt,
	}
}

func (b bearerAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+b.bearerToken)
	return b.rt.RoundTrip(req)
}
