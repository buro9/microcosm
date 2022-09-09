package transform

import (
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// NaturalTime converts a time into a moment relative to now in human terms.
// i.e. `12 secondds ago` or `3 days from now`
func NaturalTime(d time.Time) string {
	return humanize.Time(d)
}

// RFCTime converts a time into the RFC3339 representation of that time for use
// in a HTML5 time element.
func RFCTime(d time.Time) string {
	return d.UTC().Format(time.RFC3339)
}

// MsToSeconds converts a time in milliseconds to a representation in seconds.
// i.e. 147 milliseconds = 0.147
func MsToSeconds(value interface{}) string {
	var ms float64
	switch v := value.(type) {
	case float32:
		ms = float64(v)
	case float64:
		ms = v
	case int:
		ms = float64(v)
	case int32:
		ms = float64(v)
	case int64:
		ms = float64(v)
	default:
		return ""
	}

	if ms == 0 {
		return ""
	}

	return strconv.FormatFloat(ms/1000, 'f', 3, 64)
}
