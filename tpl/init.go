package tpl

import (
	"html/template"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/hugo/helpers"
)

// FuncMap is a public export of the Hugo FuncMap
var FuncMap template.FuncMap

func init() {
	FuncMap = template.FuncMap{
		"absURL": absURL,
		"absLangURL": func(i interface{}) (template.HTML, error) {
			s, err := cast.ToStringE(i)
			if err != nil {
				return "", err
			}
			return template.HTML(helpers.CurrentPathSpec().AbsURL(s, true)), nil
		},
		"add":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '+') },
		"after":        after,
		"base64Decode": base64Decode,
		"base64Encode": base64Encode,
		"chomp":        chomp,
		"countrunes":   countRunes,
		"countwords":   countWords,
		"default":      dfault,
		"dateFormat":   dateFormat,
		"delimit":      delimit,
		"dict":         dictionary,
		"div":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '/') },
		"echoParam":    returnWhenSet,
		"emojify":      emojify,
		"eq":           eq,
		"findRE":       findRE,
		"first":        first,
		"ge":           ge,
		"getCSV":       getCSV,
		"getJSON":      getJSON,
		"getenv":       getenv,
		"gt":           gt,
		"hasPrefix":    hasPrefix,
		"highlight":    highlight,
		"htmlEscape":   htmlEscape,
		"htmlUnescape": htmlUnescape,
		"humanize":     humanize,
		"imageConfig":  imageConfig,
		"in":           in,
		"index":        index,
		"int":          func(v interface{}) (int, error) { return cast.ToIntE(v) },
		"intersect":    intersect,
		"isSet":        isSet,
		"isset":        isSet,
		"jsonify":      jsonify,
		"last":         last,
		"le":           le,
		"lower":        lower,
		"lt":           lt,
		"markdownify":  markdownify,
		"md5":          md5,
		"mod":          mod,
		"modBool":      modBool,
		"mul":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '*') },
		"ne":           ne,
		"now":          func() time.Time { return time.Now() },
		"plainify":     plainify,
		"pluralize":    pluralize,
		"querify":      querify,
		"readDir":      readDirFromWorkingDir,
		"readFile":     readFileFromWorkingDir,
		"ref":          ref,
		"relURL":       relURL,
		"relLangURL": func(i interface{}) (template.HTML, error) {
			s, err := cast.ToStringE(i)
			if err != nil {
				return "", err
			}
			return template.HTML(helpers.CurrentPathSpec().RelURL(s, true)), nil
		},
		"relref":       relRef,
		"replace":      replace,
		"replaceRE":    replaceRE,
		"safeCSS":      safeCSS,
		"safeHTML":     safeHTML,
		"safeHTMLAttr": safeHTMLAttr,
		"safeJS":       safeJS,
		"safeURL":      safeURL,
		"sanitizeURL":  helpers.SanitizeURL,
		"sanitizeurl":  helpers.SanitizeURL,
		"seq":          helpers.Seq,
		"sha1":         sha1,
		"sha256":       sha256,
		"shuffle":      shuffle,
		"singularize":  singularize,
		"slice":        slice,
		"slicestr":     slicestr,
		"sort":         sortSeq,
		"split":        split,
		"string":       func(v interface{}) (string, error) { return cast.ToStringE(v) },
		"sub":          func(a, b interface{}) (interface{}, error) { return helpers.DoArithmetic(a, b, '-') },
		"substr":       substr,
		"title":        title,
		"time":         asTime,
		"trim":         trim,
		"truncate":     truncate,
		"upper":        upper,
		"where":        where,
		"i18n":         i18nTranslate,
		"T":            i18nTranslate,
	}
}
