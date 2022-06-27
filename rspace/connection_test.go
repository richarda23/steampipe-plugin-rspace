package rspace

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	BASE_URL := "https://myrspace.somewhere.com"
	EXPECTED_URL, _ := url.Parse("https://myrspace.somewhere.com/api/v1")
	parsed, _ := parseUrl(BASE_URL)
	assert.Equal(t, EXPECTED_URL, parsed)

	parsed, _ = parseUrl(BASE_URL + "/")
	assert.Equal(t, EXPECTED_URL, parsed)

	parsed, _ = parseUrl(BASE_URL + "/api/v1")
	assert.Equal(t, EXPECTED_URL, parsed)
}
