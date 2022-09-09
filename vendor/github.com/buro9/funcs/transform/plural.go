package transform

import "github.com/dustin/go-humanize/english"

// Plural takes a number and a singular word and pluralises using basic English
// rules. If the output is not as expected, provide the plural word too.
//
// e.g. Plural(17, "record", "") = 17 records
func Plural(value interface{}, singular string, plural string) string {
	var quantity int
	switch v := value.(type) {
	case float32:
		quantity = int(v)
	case float64:
		quantity = int(v)
	case int:
		quantity = v
	case int32:
		quantity = int(v)
	case int64:
		quantity = int(v)
	default:
		return ""
	}
	return english.Plural(quantity, singular, plural)
}

// PluralWord takes a number and a singular word and pluralises using basic English
// rules. If the output is not as expected, provide the plural word too.
//
// e.g. PluralWord(17, "record", "") = records
func PluralWord(value interface{}, singular string, plural string) string {
	var quantity int
	switch v := value.(type) {
	case float32:
		quantity = int(v)
	case float64:
		quantity = int(v)
	case int:
		quantity = v
	case int32:
		quantity = int(v)
	case int64:
		quantity = int(v)
	default:
		return ""
	}
	return english.PluralWord(quantity, singular, plural)
}
