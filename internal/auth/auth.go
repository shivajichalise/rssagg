package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts API key from request headers
// Example:
// Authorization: APIKey {insert api key here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found.")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header.")
	}

	if vals[0] != "APIKey" {
		return "", errors.New("malformed first part of auth header.")
	}

	return vals[1], nil
}
