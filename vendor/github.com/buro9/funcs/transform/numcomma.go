package transform

import humanize "github.com/dustin/go-humanize"

// NumComma takes a number and returns it comma delimited
//
// If the value provided is not a number, an empty string is returned
func NumComma(value interface{}) string {
	switch v := value.(type) {
	case float32:
		return humanize.Commaf(float64(v))
	case float64:
		return humanize.Commaf(v)
	case int:
		return humanize.Comma(int64(v))
	case int32:
		return humanize.Comma(int64(v))
	case int64:
		return humanize.Comma(v)
	default:
		return ""
	}
}
