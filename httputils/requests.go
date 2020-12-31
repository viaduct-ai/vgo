package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

const (
	// ContentTypeJSON is a constant for JSON Header Type
	ContentTypeJSON = "application/json"

	oneMB = 1048576
)

var (
	// ContentType is a constant for the Content-Type header
	ContentType = http.CanonicalHeaderKey("Content-Type")
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func (mr *malformedRequest) Status() int {
	return mr.status
}

func (mr *malformedRequest) Message() string {
	return mr.msg
}

func (mr *malformedRequest) Code() string {
	return "bad_request"
}

// DecodeJSONBody strictly decodes a JSON request body into a given interface.
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get(ContentType) != "" {
		value, _ := header.ParseValueAndParams(r.Header, ContentType)
		if value != ContentTypeJSON {
			msg := fmt.Sprintf("%s header is not %s", ContentType, ContentTypeJSON)
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, oneMB)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

			// There is an open issue at https://github.com/golang/go/issues/29035 regarding
			// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

			// Again there is an open issue regarding turning this into a sentinel
			// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
