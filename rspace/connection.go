package rspace

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

func getenv(envname string) string {
	return os.Getenv(envname)
}

func parseUrl(urlStr string) (*url.URL, error) {
	var canonicalUrl string = urlStr
	if !strings.HasSuffix(urlStr, "/api/v1") {
		if !strings.HasSuffix(urlStr, "/") {
			canonicalUrl = urlStr + "/"
		}
		canonicalUrl = canonicalUrl + "api/v1"
	}
	apiUrl, err := url.Parse(canonicalUrl)
	return apiUrl, err
}

func connect(ctx context.Context) (*rspace.RsWebClient, error) {
	url, e := parseUrl(getenv(BASE_URL_ENV_NAME))
	if e != nil {
		return nil, e
	}
	apikey := getenv(APIKEY_ENV_NAME)
	if apikey == "" || url == nil {
		return nil, errors.New("API key and/ or URL is not defined")
	}
	webClient := rspace.NewWebClientCustomTimeout(url, apikey, 30)
	return webClient, nil
}
