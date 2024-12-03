package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")

	if apiKey == "" {
		return "", errors.New("no api key provided")
	}

	keys := strings.Split(apiKey, " ")

	if len(keys) != 2 {
		return "", errors.New("invalid api key")
	}

	if keys[0] != "ApiKey" {
		return "", errors.New("invalid api key")
	}

	return keys[1], nil
}
