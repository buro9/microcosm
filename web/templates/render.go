package templates

import (
	"fmt"
	"net/http"

	"github.com/oxtoacart/bpool"
)

var (
	// We buffer template executions so that we can catch errors
	bufpool *bpool.BufferPool
)

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
	tmpl := templates[name].Lookup(name + ".html.tmpl")
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
