package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extract apikey from headres of http request
// Authorization: ApiKey  <key>
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("could not find auth info")
	}

	values := strings.Split(val, " ")

	if len(values) != 2 {
		return "", errors.New("Malformed header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("Malformed first part of header")
	}

	return values[1], nil

}
