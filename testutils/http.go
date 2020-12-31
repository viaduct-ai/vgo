package testutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServer allows convenience testing methods on-top of the standard httptest,Server
type TestServer struct {
	*httptest.Server
}

// NewTestServer returns a new a server for testing
func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewServer(h)
	return &TestServer{ts}
}

// Get request to the test server
func (ts *TestServer) Get(t *testing.T, path string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

// Post request to the test server
func (ts *TestServer) Post(t *testing.T, path string, reqBody []byte) (int, http.Header, []byte) {
	rs, err := ts.Client().Post(ts.URL+path, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
