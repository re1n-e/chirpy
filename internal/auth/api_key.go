package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	authString := header.Get("Authorization")
	if authString == "" {
		return "", errors.New("No Authorizatin header found")
	}
	resp := strings.Fields(authString)
	if len(resp) != 2 || resp[0] != "ApiKey" {
		return "", errors.New("Invalid header format")
	}
	return resp[1], nil
}
