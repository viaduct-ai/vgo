package jwtutils

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	jwtRequest "github.com/golang-jwt/jwt/request"
)

// jwtParser for strictly parsing, not validating, a JWT token.
var jwtParser = jwt.Parser{
	UseJSONNumber:        true,
	SkipClaimsValidation: true,
}

// ParseUnverifiedTokenClaimsFromRequest parses a *http.Request for a token in the Authorization header and returns the claims if the token exists and is valid. This method does NOT verify the Authorization token. It assumes token validation has happened upstream.
func ParseUnverifiedTokenClaimsFromRequest(r *http.Request) (map[string]interface{}, error) {
	rawToken, err := jwtRequest.AuthorizationHeaderExtractor.ExtractToken(r)

	if err != nil {
		return map[string]interface{}{}, err
	}

	claims := jwt.MapClaims{}

	_, _, err = jwtParser.ParseUnverified(rawToken, claims)

	if err != nil {
		return map[string]interface{}{}, err
	}

	return map[string]interface{}(claims), nil
}
