package middlewares_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/viaduct-ai/vgo/ginutils/middlewares"
)

var envoyOriginalPath = http.CanonicalHeaderKey("X-Envoy-Original-Path")

func TestEnvoyProxyMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      *http.Request
		wantPath string
		wantHost string
		wantURL  string
	}{
		{
			name: "Non-Envoy Preserves Path",
			req: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "localhost:3000",
					Path:     "/test",
					RawQuery: "q=test",
					Fragment: "test",
				},
				Host: "localhost",
			},
			wantPath: "/test",
			wantHost: "localhost:3000",
			wantURL:  "http://localhost:3000/test?q=test#test",
		},
		{
			name: "URL Path and Host Update",
			req: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "localhost:3000",
					Path:     "/v1/test",
					RawQuery: "q=test",
					Fragment: "test",
				},
				Host: "example.com",
				Header: http.Header(map[string][]string{
					envoyOriginalPath: {"/test"},
				}),
			},
			wantPath: "/test",
			wantHost: "example.com",
			wantURL:  "http://example.com/test?q=test#test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			middlewares.EnvoyProxyMiddleware(&gin.Context{Request: tt.req})

			// check the req was mutated correctly
			if tt.req.URL.Path != tt.wantPath {
				t.Errorf("want path %s. got %s", tt.wantPath, tt.req.URL.Path)
			}

			if tt.req.URL.Host != tt.wantHost {
				t.Errorf("want host %s. got %s", tt.wantHost, tt.req.URL.Host)
			}

			if tt.req.URL.String() != tt.wantURL {
				t.Errorf("want host %s. got %s", tt.wantURL, tt.req.URL.String())
			}
		})
	}
}
