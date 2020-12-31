package handlers

import (
	"net/http"

	"github.com/viaduct-ai/vgo/httputils"
)

// HealthCheck handler is a standard handler for checking if server is alive or ready
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputils.ServeJSON(w, http.StatusOK, map[string]bool{
		"alive": true,
	})
}
