package errors

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

func Render(w http.ResponseWriter, r *http.Request, status int, err error) {
	switch status {
	case http.StatusForbidden:
		// All instances of 403 should return a 403
		// https://www.rfc-editor.org/rfc/rfc6750 3.1 insufficient_scope
		// This is due to the current authentication (if present) not being
		// sufficient to render the item. In effect, we're saying for the
		// current credentials re-authenticating will no nothing.
		w.WriteHeader(http.StatusForbidden)
		data := templates.Data{
			Request: r,
			Site:    bag.GetSite(r.Context()),
			User:    bag.GetProfile(r.Context()),

			Error:      err,
			StatusCode: status,
		}
		err = templates.RenderHTML(w, "403", data)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<h1>403 Forbidden</h1>`))
			w.Write([]byte(`<p>` + err.Error() + `</p>`))
		}
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		data := templates.Data{
			Request: r,
			Site:    bag.GetSite(r.Context()),
			User:    bag.GetProfile(r.Context()),

			Error:      err,
			StatusCode: status,
		}
		err = templates.RenderHTML(w, "404", data)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<h1>404 Not Found</h1>`))
			w.Write([]byte(`<p>` + err.Error() + `</p>`))
		}
	case http.StatusInternalServerError:
		w.WriteHeader(http.StatusInternalServerError)
		data := templates.Data{
			Request: r,
			Site:    bag.GetSite(r.Context()),
			User:    bag.GetProfile(r.Context()),

			Error:      err,
			StatusCode: status,
		}
		err = templates.RenderHTML(w, "500", data)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<h1>500 Internal Server Error</h1>`))
			w.Write([]byte(`<p>` + err.Error() + `</p>`))
		}
	default:
		fmt.Printf("TODO: ERROR HANDLER NEEDED: %d %s\n", status, err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<h1>Undefined Error</h1>`))
		w.Write([]byte(err.Error()))
	}
}
