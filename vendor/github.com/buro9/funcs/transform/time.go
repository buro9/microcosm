package transform

import (
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
