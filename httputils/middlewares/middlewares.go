package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/viaduct-ai/vgo/jwtutils"
	"github.com/viaduct-ai/vgo/log"
)

// Custom http middleware
// Other common middleware can be found at
// https://github.com/gorilla/handlers
var (
	userAgent = http.CanonicalHeaderKey("User-Agent")
	referer   = http.CanonicalHeaderKey("Referer")

	// X-Envoy-Original-Path
	envoyOriginalPath = http.CanonicalHeaderKey("X-Envoy-Original-Path")
)

// TODO: Decide about including gorilla proxy headers

// LoggingMiddleware logs all incoming requests and response headers using the internal logger
func LoggingMiddleware(l log.Logger, next http.Handler) http.Handler {
	logRequest := func(w http.ResponseWriter, r *http.Request) {

		// try parse claims from request, ignore errors
		claims, _ := jwtutils.ParseUnverifiedTokenClaimsFromRequest(r)

		// read, close, re-create
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		rData := map[string]interface{}{}

		if err == nil {
			json.Unmarshal(body, &rData)

			// Omit password from logs
			// Need to find a more generic approach for any sensitive fields
			delete(rData, "password")
		}

		l.With(
			"method", r.Method,
			"host", r.Host,
			"url", r.URL.String(),
			"endpoint", r.URL.Path,
			"scheme", r.URL.Scheme,
			"ip", r.RemoteAddr,
			"headers", r.Header,
			"query", r.URL.Query(),
			"body", rData,
			"auth", claims,
		).Info("request")
	}

	logHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(w, r)
		next.ServeHTTP(w, r)
		// TODO: Log Response Code and Metrics
		// This is harder than expected, consider using https://github.com/felixge/httpsnoop
		// logResponse(w, r)
	})

	return logHandler
}

// EnvoyProxyMiddleware modifies the request with envoy proxy specific headers
//
// Ideally, we should not add additional header information here.
// This should be through envoy header manipulation
// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#custom-request-response-headers
func EnvoyProxyMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if path := r.Header.Get(envoyOriginalPath); path != "" {
			orgURL, err := url.Parse(path)

			if err == nil {
				r.URL.Path = orgURL.Path
			}

			// https://stackoverflow.com/questions/42921567/what-is-the-difference-between-host-and-url-host-for-golang-http-request
			// The request went through an envoy proxy, so update the URL host to match the original reqest host
			r.URL.Host = r.Host
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
