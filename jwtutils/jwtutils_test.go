package jwtutils_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	jwtRequest "github.com/golang-jwt/jwt/request"

	"github.com/viaduct-ai/vgo/jwtutils"
)

var (
	testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func TestParseUnverifiedTokenClaimsFromRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		headers    http.Header
		wantClaims map[string]interface{}
		wantErr    error
	}{
		{
			name:       "Missing Header",
			wantErr:    jwtRequest.ErrNoTokenInRequest,
			wantClaims: map[string]interface{}{},
		},
		{
			name: "Malformed Header",
			headers: http.Header(map[string][]string{
				"Authorization": {"Broken eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			}),
			wantErr:    base64.CorruptInputError(6),
			wantClaims: map[string]interface{}{},
		}, {
			name: "Malformed JWT",
			headers: http.Header(map[string][]string{
				"Authorization": {"Bearer eyJhbGciOiJzI1NiIsInR5cCI6IkpXVCJ9.eyJzTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			}),
			wantErr:    errors.New("invalid character"),
			wantClaims: map[string]interface{}{},
		}, {
			name: "Valid",
			headers: http.Header(map[string][]string{
				"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			}),
			wantClaims: map[string]interface{}{
				"sub":  "1234567890",
				"name": "John Doe",
				"iat":  json.Number("1516239022"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				Header: tt.headers,
			}

			claims, err := jwtutils.ParseUnverifiedTokenClaimsFromRequest(r)

			if err != tt.wantErr && !errors.Is(err, tt.wantErr) && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("want %v. got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(claims, tt.wantClaims) {
				t.Errorf("want %v. got %v", tt.wantClaims, claims)
			}
		})
	}
}
