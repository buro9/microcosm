// Package sanitize provides funcs that will allow untrusted content into a
// template and controls the escaping.
package sanitize

import (
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

var (
	textPolicy     = bluemonday.StripTagsPolicy()
	htmlPolicy     = bluemonday.UGCPolicy()
	initHTMLPolicy bool
)

// ASCIISpace ...
var ASCIISpace = rune(` `[0])

// StripTags allows everything, allows tabs, newlines, ASCII space but
// non-conforming whitespace, control chars, and also prevents shouting
func StripTags(s string) string {
	return textPolicy.Sanitize(stripChars(s, true, true, true, true))
}

// StripTagsSentence strips all HTML tags, allows ASCII space
func StripTagsSentence(s string) string {
	return textPolicy.Sanitize(stripChars(s, true, true, false, true))
}

// SanitiseHTML sanitizes HTML
// Leaving a safe set of HTML intact that is not going to pose an XSS risk
func Sanitize(s string) string {
	if !initHTMLPolicy {
		htmlPolicy.RequireNoFollowOnLinks(false)
		htmlPolicy.RequireNoFollowOnFullyQualifiedLinks(true)
		htmlPolicy.AddTargetBlankToFullyQualifiedLinks(true)
		initHTMLPolicy = true
	}

	return htmlPolicy.Sanitize(s)
}

// stripChars will remove unicode characters according to the instructions given
// to it. This is strictly speaking a whitelist, it's less "strip", more "allow".
//
// With SanitiseText this should run before bluemonday (HTML sanitiser)
//
// With SanitiseHTML this should run before blackfriday (Markdown generator)
func stripChars(
	in string,
	allowPrint bool,
	allowASCIISpace bool,
	allowFormattingSpace bool,
	preventShouting bool,
) string {
	tmp := []rune{}

	isShouting := true

	for _, runeValue := range in {
		var ret bool
		// IsPrint covers anything actually printable plus ASCII space
		if allowPrint && unicode.IsPrint(runeValue) {
			if !allowASCIISpace && unicode.IsSpace(runeValue) {
				continue
			}

			if isShouting && preventShouting && unicode.IsLower(runeValue) {
				isShouting = false
			}

			ret = true
		}

		// Only ASCII space
		if allowASCIISpace && runeValue == ASCIISpace {
			ret = true
		}

		// IsSpace covers tabs, newlines, ASCII space,etc
		if allowFormattingSpace && unicode.IsSpace(runeValue) {
			ret = true
		}

		if ret {
			tmp = append(tmp, runeValue)
		}
	}

	if isShouting && preventShouting {
		tmp2 := []rune{}
		for _, runeValue := range string(tmp) {
			tmp2 = append(tmp2, unicode.ToLower(runeValue))
		}
		tmp = tmp2
	}

	return string(tmp)
}
