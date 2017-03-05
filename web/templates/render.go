package templates

import (
	"fmt"
	"net/http"

	"github.com/oxtoacart/bpool"

	"github.com/buro9/microcosm/models"
)

var (
	// We buffer template executions so that we can catch errors
	bufpool *bpool.BufferPool
)

// Data is the data that can be provided to a template.
//
// This is normalised into this one struct to ensure consistency across all
// templates, though very obviously not all templates require all fields and
// most of the time very few fields are filled in, typically an anonymous user
// will only have Site and whatever fields are relavent for a page shown, and a
// signed-in user will have Site and User along with whatever fields are
// relevant for the current page.
type Data struct {
	// Every request
	Request    *http.Request
	Site       *models.Site
	Section    string
	Query      *models.SearchQuery
	Pagination *models.Pagination

	// If signed-in
	User *models.Profile

	// Depending on context, templates may expect the applicable one to be
	// filled in
	Array         *models.Array
	Microcosm     *models.Microcosm
	SearchResults *models.SearchResults
}

// RenderHTML is a wrapper around template.ExecuteTemplate.
//
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to
// catch any errors resulting from populating the template.
func RenderHTML(
	w http.ResponseWriter,
	name string,
	data Data,
) error {
	// Ensure the template exists in the map.
	//
	// The map exists in templates.go and is populated at init from the
	// definitions held in definitions.go
	tmpl := templates[name].Lookup(name + ".tmpl")
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
