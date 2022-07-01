package rspace

import "strings"

// date_from_iso8601 extracts the date part of an ISO8601 timestamp
func date_from_iso8601(iso8601 string) string {
	parts := strings.Split(iso8601, "T")
	return parts[0]
}
