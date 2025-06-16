package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestOptions struct {
	Method  string
	Path    string
	Body    io.Reader
	Headers map[string]string
}

func MakeRequest(t *testing.T, handler http.Handler, opts RequestOptions) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(opts.Method, opts.Path, opts.Body)

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	handler.ServeHTTP(rr, req)

	return rr
}
