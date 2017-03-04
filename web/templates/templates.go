package templates

// Templates is a slice of template definitions that will be loaded and compiled
// during the initialisation of this package
var Templates []Template

// Template is a description of a HTML page as a collection of template files.
type Template struct {
	// Base represents the template file at
	// {*opts.FilesPath}/base/{template.Base}.tmpl
	//
	// This is expected to contain the HTML header, navigation and footer as
	// well as the core structure of the HTML page.
	Base string

	// Page represents the template file at
	// {*opts.FilesPath}/pages/{template.Page}.tmpl
	//
	// This is expected to be the content that populates the Base template
	Page string

	// Includes is an optional list of template files at
	// {*opts.FilesPath}/includes/{template.Includes[i]}.tmpl
	//
	// These are expected to be common blocks that multiple pages may use
	Includes []string
}

// Collate consumes any number of strings or slices of strings and constructs a
// single slice from them.
func Collate(includes ...interface{}) (out []string) {
	for _, include := range includes {
		switch v := include.(type) {
		case string:
			out = append(out, v)
		case []string:
			for _, s := range v {
				out = append(out, s)
			}
		}
	}
	return
}
