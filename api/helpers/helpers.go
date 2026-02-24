package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL, nil
	}
	return rawURL, nil
}

func RemoveDomainError(rawUrl string) bool {
	if rawUrl == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(rawUrl, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "wwww.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
