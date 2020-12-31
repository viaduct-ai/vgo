package middlewares_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/viaduct-ai/vgo/httputils/middlewares"
	"github.com/viaduct-ai/vgo/testutils"
	"golang.org/x/net/context"
)

var envoyOriginalPath = http.CanonicalHeaderKey("X-Envoy-Original-Path")

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	// no-op
}

func TestEnvoyProxyMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(dummyHandler)
	middleware := middlewares.EnvoyProxyMiddleware(handler)

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

			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, tt.req)

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

func TestLoggingMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(dummyHandler)

	var testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	baseReq := httptest.NewRequest(http.MethodPost, "https://api.test.com/v1/token", nil)

	authReq := baseReq.Clone(context.Background())
	authReq.Header = http.Header(map[string][]string{
		"Authorization": []string{fmt.Sprintf("Bearer %s", testToken)},
	})

	loginReq := authReq.Clone(context.Background())
	loginReq.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{
		"username": "test",
		"password": "test"
	}`)))

	// validate password is not logged
	tests := []struct {
		name        string
		req         *http.Request
		wantLogs    []interface{}
		wantContext map[string]interface{}
	}{
		{
			name:     "Valid",
			req:      baseReq,
			wantLogs: []interface{}{"request"},
			wantContext: map[string]interface{}{
				"method":   baseReq.Method,
				"host":     baseReq.Host,
				"url":      baseReq.URL.String(),
				"endpoint": baseReq.URL.Path,
				"scheme":   baseReq.URL.Scheme,
				"ip":       baseReq.RemoteAddr,
				"headers":  baseReq.Header,
				"query":    baseReq.URL.Query(),
				"body":     map[string]interface{}{},
				"auth":     map[string]interface{}{},
			},
		},
		{
			name:     "Authentication",
			req:      authReq,
			wantLogs: []interface{}{"request"},
			wantContext: map[string]interface{}{
				"method":   authReq.Method,
				"host":     authReq.Host,
				"url":      authReq.URL.String(),
				"endpoint": authReq.URL.Path,
				"scheme":   authReq.URL.Scheme,
				"ip":       authReq.RemoteAddr,
				"headers":  authReq.Header,
				"query":    authReq.URL.Query(),
				"body":     map[string]interface{}{},
				"auth": map[string]interface{}{
					"iat":  json.Number("1516239022"),
					"name": "John Doe",
					"sub":  "1234567890",
				},
			},
		},
		{
			name:     "Body Logged and Password Ignored",
			req:      loginReq,
			wantLogs: []interface{}{"request"},
			wantContext: map[string]interface{}{
				"method":   loginReq.Method,
				"host":     loginReq.Host,
				"url":      loginReq.URL.String(),
				"endpoint": loginReq.URL.Path,
				"scheme":   loginReq.URL.Scheme,
				"ip":       loginReq.RemoteAddr,
				"headers":  loginReq.Header,
				"query":    loginReq.URL.Query(),
				"body": map[string]interface{}{
					"username": "test",
				},
				"auth": map[string]interface{}{
					"iat":  json.Number("1516239022"),
					"name": "John Doe",
					"sub":  "1234567890",
				},
			},
		},
		// {name: "Password Not Logged"},
	}

	// Create  test logger to capture arguments
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := testutils.NewTestLogger()

			middleware := middlewares.LoggingMiddleware(logger, handler)

			rr := httptest.NewRecorder()

			// read, close, re-create
			wantBody, _ := ioutil.ReadAll(tt.req.Body)
			tt.req.Body.Close()
			tt.req.Body = ioutil.NopCloser(bytes.NewBuffer(wantBody))

			middleware.ServeHTTP(rr, tt.req)

			body, err := ioutil.ReadAll(tt.req.Body)
			if err != nil {
				t.Errorf("error reading body after middleware")
			}

			if bytes.Compare(wantBody, body) != 0 {
				t.Errorf("want body %s. got %s", wantBody, body)
			}

			// validate info logs
			if !reflect.DeepEqual(tt.wantLogs, logger.InfoLogs) {
				t.Errorf("want info logs %v. got %v", tt.wantLogs, logger.InfoLogs)
			}

			// validate context
			if !reflect.DeepEqual(tt.wantContext, logger.Context) {
				t.Errorf("want context %v. got %v", tt.wantContext, logger.Context)
			}
		})
	}
}
