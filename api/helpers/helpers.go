package helpers

import (
	"net/url"
	"os"
	"strings"
)

func EnforceHTTP(rawUrl string) (string, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	if url.Scheme != "http" {
		return "https://" + rawUrl, nil
	}
	return rawUrl, nil
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
