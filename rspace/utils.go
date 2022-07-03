package rspace

import (
	"regexp"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

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

type dateQuery struct {
	from, to time.Time
}

func (d dateQuery) String() string {
	rc := ""
	if !d.from.IsZero() {
		rc = d.from.Format("2006-01-02")
	}
	rc = rc + ";"
	if !d.to.IsZero() {
		rc = rc + d.to.Format("2006-01-02")
	}
	return rc
}

func getDateQueryFromQual(fieldName string, quals plugin.KeyColumnQualMap) dateQuery {
	dateQ := dateQuery{}
	for _, q := range quals[fieldName].Quals {
		val := q.Value.GetTimestampValue()
		switch q.Operator {
		case ">", ">=":
			dateQ.from = val.AsTime()
		case "<", "<=":
			dateQ.to = val.AsTime()
		}
	}
	return dateQ
}
