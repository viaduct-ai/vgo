package httputils

import (
	"encoding/json"
	"errors"

	"net/http"
)

var internalError = APIErrorResponse{
	Message: "an internal error has occurred",
	Code:    "internal",
}

// APIError is an interface for creating an API error message.
// This allows users to implement their own custom error struct for their applications
// while still making them compatible with the APIError interface. The interface also allows users// to type assert errors without requiring them to import the APIError interface
// Inspired by https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
type APIError interface {
	Error() string
	Status() int
	Message() string
	Code() string
}

// APIErrorResponse is a simple API error response format.
// In the future, we can consider a more comprehensive and standard error format,
// such as application/problem+json
type APIErrorResponse struct {
	Message string
	Code    string
}

// ServeError serves an APIErrorResponse.
// If the err implements the APIError interface, its content will be used in the response.
// Else it will serve an internal error response.
func ServeError(w http.ResponseWriter, err error) {
	var apiError APIError
	if errors.As(err, &apiError) {
		resp := APIErrorResponse{
			Message: apiError.Message(),
			Code:    apiError.Code(),
		}

		ServeJSON(w, apiError.Status(), resp)
		return
	}

	ServeJSON(w, http.StatusInternalServerError, internalError)
}

// ServeJSON is serves a formatted, JSON response the user
func ServeJSON(w http.ResponseWriter, status int, body interface{}) {
	resp, err := json.MarshalIndent(body, "", "\t")

	if err != nil {
		ServeError(w, err)
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(status)
	w.Write(resp)
}
