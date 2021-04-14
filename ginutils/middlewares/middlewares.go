package middlewares

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Custom gin middleware
var (
	// X-Envoy-Original-Path
	envoyOriginalPath = http.CanonicalHeaderKey("X-Envoy-Original-Path")
)

// EnvoyProxyMiddleware modifies the request with envoy proxy specific headers
//
// Ideally, we should not add additional header information here.
// This should be through envoy header manipulation
// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#custom-request-response-headers
func EnvoyProxyMiddleware(c *gin.Context) {
	if path := c.Request.Header.Get(envoyOriginalPath); path != "" {
		orgURL, err := url.Parse(path)

		if err == nil {
			c.Request.URL.Path = orgURL.Path
		}

		// https://stackoverflow.com/questions/42921567/what-is-the-difference-between-host-and-url-host-for-golang-http-request
		// The request went through an envoy proxy, so update the URL host to match the original reqest host
		c.Request.URL.Host = c.Request.Host
	}

	c.Next() // Pass on to the next-in-chain
}
