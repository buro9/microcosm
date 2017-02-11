// Package safe provides funcs that will allow trusted content into a template without
// being escaped.
//
// Portions Copyright The Hugo Authors and covered by both an MIT license for
// the original code, and an Apache license for later modifications.
// https://github.com/spf13/hugo/blob/master/tpl/template_funcs.go
package safe

import (
	"html/template"

	"github.com/spf13/cast"
)

// CSS returns a given string as html/template CSS content
//
// Use of this func presents a security risk: the encapsulated content should
// come from a trusted source, as it will be included verbatim in the template
// output
func CSS(a interface{}) (template.CSS, error) {
	s, err := cast.ToStringE(a)
	return template.CSS(s), err
}

// HTML returns a given string as html/template HTML content
//
// Use of this func presents a security risk: the encapsulated content should
// come from a trusted source, as it will be included verbatim in the template
// output
func HTML(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	return template.HTML(s), err
}

// HTMLAttr returns a given string as html/template HTMLAttr content
//
// Use of this func presents a security risk: the encapsulated content should
// come from a trusted source, as it will be included verbatim in the template
// output
func HTMLAttr(a interface{}) (template.HTMLAttr, error) {
	s, err := cast.ToStringE(a)
	return template.HTMLAttr(s), err
}

// JS returns the given string as a html/template JS content
//
// Use of this func presents a security risk: the encapsulated content should
// come from a trusted source, as it will be included verbatim in the template
// output
func JS(a interface{}) (template.JS, error) {
	s, err := cast.ToStringE(a)
	return template.JS(s), err
}

// URL returns a given string as html/template URL content
//
// Use of this func presents a security risk: the encapsulated content should
// come from a trusted source, as it will be included verbatim in the template
// output
func URL(a interface{}) (template.URL, error) {
	s, err := cast.ToStringE(a)
	return template.URL(s), err
}
