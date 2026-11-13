package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key := strings.Fields(headers.Get("Authorization"))
	if len(key) != 2 || !strings.EqualFold(key[0], "ApiKey") {
		return "", fmt.Errorf("Couldn't generate api key try again later")
	}
	return key[1], nil
}
