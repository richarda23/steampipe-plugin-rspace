package rspace

import (
	"regexp"
	"strconv"
	"strings"
)

// date_from_iso8601 extracts the date part of an ISO8601 timestamp
func date_from_iso8601(iso8601 string) string {
	parts := strings.Split(iso8601, "T")
	return parts[0]
}

func isGlobalId(idStr string) bool {
	if idStr == "" {
		return false
	}
	ok, _ := regexp.MatchString(`^[A-Z]{2}\d+`, idStr)
	return ok
}

func getIdFromGlobalId(idStr string) (int, error) {
	if isGlobalId(idStr) {
		return strconv.Atoi(idStr[2:])
	} else {
		return strconv.Atoi(idStr)
	}
}
