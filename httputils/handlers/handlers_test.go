package handlers_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/viaduct-ai/vgo/httputils/handlers"
	"github.com/viaduct-ai/vgo/testutils"
)

func TestHealtCheck(t *testing.T) {
	path := "/health"
	mux := http.NewServeMux()
	mux.HandleFunc(path, handlers.HealthCheck)

	ts := testutils.NewTestServer(t, mux)

	code, _, body := ts.Get(t, path)

	if code != http.StatusOK {
		t.Errorf("want status code %d. got %d", http.StatusOK, code)
	}

	wantBody := map[string]bool{
		"alive": true,
	}

	var gotBody map[string]bool
	json.Unmarshal(body, &gotBody)

	if !reflect.DeepEqual(wantBody, gotBody) {
		t.Errorf("want %v. got %v", wantBody, gotBody)
	}
}
