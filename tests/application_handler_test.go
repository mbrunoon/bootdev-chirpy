package tests

import (
	"net/http"
	"testing"

	"github.com/mbrunoon/bootdev-chirpy/internal/application"
	"github.com/stretchr/testify/assert"
)

func TestHealthzHandler(t *testing.T) {
	handler := http.HandlerFunc(application.HealthzHandler)

	opts := RequestOptions{
		Method: http.MethodGet,
		Path:   "/healthz",
	}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /healthz 200",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "OK"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := MakeRequest(t, handler, opts)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, rr.Body.String())
		})
	}
}
