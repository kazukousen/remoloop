package server

import (
	"net/http"
	"strings"

	"github.com/go-kit/kit/log/level"

	"github.com/kazukousen/remoloop/pkg/remoloop/api"
)

func (s server) rewriteRootPath(path string, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = r.URL.Path[len(path):]
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	level.Debug(s.logger).Log("msg", "get http request", "method", r.Method, "path", r.URL.Path)
	if r.URL.Path == "/ready" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		resource := r.URL.Path[len("/api"):]
		s.client.Get(r.Context(), api.Resource(resource), w)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("this page not found"))
	return
}
