package httputils_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/viaduct-ai/vgo/httputils"
)

func TestDecodeJSONBody(t *testing.T) {
	t.Parallel()

	type test struct {
		Test string `json:"test,omitempty"`
	}

	tests := []struct {
		name         string
		header       http.Header
		body         []byte
		decodeStruct test
		wantBody     test
		wantErr      error
	}{
		{
			name: "Valid",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"test": "test"
			}`),
			decodeStruct: test{},
			wantBody: test{
				Test: "test",
			},
		},
		{
			name: "Wrong Content-Type",
			header: http.Header(map[string][]string{
				httputils.ContentType: {"test"},
			}),
			body: []byte(`{
				"test": "test"
			}`),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("Content-Type header is not"),
		},
		{
			name: "Unknown Field",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"unknown": "unknown"
			}`),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("unknown field"),
		},
		{
			name: "Bad JSON",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"test": true
			`),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("badly-formed"),
		},
		{
			name: "Bad JSON Syntax",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"test": "test",
			}`),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("badly-formed"),
		},
		{
			name: "Invalid Value",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"test": true
			}`),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("invalid value"),
		},
		{
			name: "Empty JSON",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			decodeStruct: test{},
			wantBody:     test{},
			wantErr:      errors.New("must not be empty"),
		},
		{
			name: "Multiple JSON",
			header: http.Header(map[string][]string{
				httputils.ContentType: {httputils.ContentTypeJSON},
			}),
			body: []byte(`{
				"test": "test"
			}, { "test": "test" }`),
			decodeStruct: test{},
			wantBody: test{
				Test: "test",
			},
			wantErr: errors.New("single JSON object"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r := &http.Request{
				Header: tt.header,
				Body:   ioutil.NopCloser(bytes.NewBuffer(tt.body)),
			}

			err := httputils.DecodeJSONBody(rr, r, &tt.decodeStruct)

			if err != tt.wantErr && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("want %v. got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.decodeStruct, tt.wantBody) {
				t.Errorf("want %+v (%T). got %+v (%T)", tt.wantBody, tt.wantBody, tt.decodeStruct, tt.decodeStruct)
			}
		})
	}
}
