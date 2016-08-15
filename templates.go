package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

var (
	templates map[string]*template.Template

	// We buffer template executions so that we can catch errors
	bufpool *bpool.BufferPool
)

// templateData is the data that can be provided to a template.
//
// This is normalised into this one struct to ensure consistency across all
// templates, though very obviously not all templates require all fields and
// most of the time very few fields are filled in, typically an anonymous user
// will only have Site and whatever fields are relavent for a page shown, and a
// signed-in user will have Site and User along with whatever fields are
// relevant for the current page.
type templateData struct {
	// Every request
	Request *http.Request
	Site    *Site
	Section string
	Query   *SearchQuery

	// If signed-in
	User *Profile

	// Depending on context, templates will expect the applicable one to be
	// filled in
	Microcosm *Microcosm
}

// Template returns the full path to a template for a given templates' name
func Template(name string) string {
	return fmt.Sprintf("%s/templates/includes/%s.tmpl", *filesPath, name)
}

// ParseTemplates loads templates on program initialisation, and is expected to
// be called by the main() func
func ParseTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layouts, err := filepath.Glob(*filesPath + "/templates/base/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	defined, err := filepath.Glob(*filesPath + "/templates/defined/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	pages, err := filepath.Glob(*filesPath + "/templates/pages/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our directories
	for _, layout := range layouts {
		files := append(append(pages, defined...), layout)

		templates[filepath.Base(layout)] =
			template.Must(
				template.New(
					filepath.Base(layout),
				).Funcs(
					funcMap(),
				).ParseFiles(
					files...,
				),
			)
	}
}

// renderHTMLTemplate is a wrapper around template.ExecuteTemplate.
//
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to
// catch any errors resulting from populating the template.
func renderHTMLTemplate(
	w http.ResponseWriter,
	name string,
	data templateData,
) error {
	// Ensure the template exists in the map.
	tmpl := templates["base.tmpl"].Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("the template named '%s' does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
