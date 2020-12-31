package httputils_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/viaduct-ai/vgo/httputils"
)

func TestServeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		body       interface{}
		wantStatus int
		wantBody   interface{}
	}{
		{
			name:   "Valid",
			status: http.StatusOK,
			body: map[string]string{
				"test": "test",
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"test": "test",
			},
		},
		{
			name:       "Empty Body",
			status:     http.StatusCreated,
			body:       map[string]string{},
			wantStatus: http.StatusCreated,
			wantBody:   map[string]string{},
		},
		{
			name:       "Nil Body",
			status:     http.StatusCreated,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Invalid Body",
			status:     http.StatusCreated,
			body:       func() {},
			wantStatus: http.StatusInternalServerError,
			wantBody: httputils.APIErrorResponse{
				Message: "an internal error has occurred",
				Code:    "internal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			httputils.ServeJSON(rr, tt.status, tt.body)

			resp := rr.Result()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("want status code %d. got %d", tt.wantStatus, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("error reading response body %v", err)
			}

			wantBody, err := json.MarshalIndent(tt.wantBody, "", "\t")

			if err != nil {
				t.Fatalf("error marshalling tt.wantBody: %v", err)
			}

			if string(body) != string(wantBody) {
				t.Errorf("want body %s. got %s", wantBody, body)
			}
		})
	}
}

type testAPIError struct {
}

func (e testAPIError) Error() string {
	return "test"
}

func (e testAPIError) Message() string {
	return "test"
}

func (e testAPIError) Code() string {
	return "test"
}

func (e testAPIError) Status() int {
	return http.StatusTeapot
}

func TestServeError(t *testing.T) {
	t.Parallel()

	testError := testAPIError{}

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   httputils.APIErrorResponse
	}{
		{
			name:       "API Error",
			err:        testError,
			wantStatus: testError.Status(),
			wantBody: httputils.APIErrorResponse{
				Message: testError.Message(),
				Code:    testError.Code(),
			},
		},
		{
			name:       "Unknown",
			err:        errors.New("unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantBody: httputils.APIErrorResponse{
				Message: "an internal error has occurred",
				Code:    "internal",
			},
		},
		{
			name:       "Nil",
			wantStatus: http.StatusInternalServerError,
			wantBody: httputils.APIErrorResponse{
				Message: "an internal error has occurred",
				Code:    "internal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			httputils.ServeError(rr, tt.err)

			resp := rr.Result()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("want status code %d. got %d", tt.wantStatus, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("error reading response body %v", err)
			}

			wantBody, err := json.MarshalIndent(tt.wantBody, "", "\t")

			if err != nil {
				t.Fatalf("error marshalling tt.wantBody: %v", err)
			}

			if string(body) != string(wantBody) {
				t.Errorf("want body %s. got %s", wantBody, body)
			}
		})
	}
}
