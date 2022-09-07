package controllers

import (
	"net/http"
	"net/url"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates"
)

// UpdatesGet will return the updates page
func UpdatesGet(w http.ResponseWriter, r *http.Request) {
	// user := bag.GetProfile(r.Context())
	// if user == nil {
	// 	w.Write([]byte("no permission"))
	// 	return
	// }

	q := url.Values{}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	updatesResults, status, err := api.GetUpdates(r.Context(), q)
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `updates`,
		Pagination: models.ParsePagination(updatesResults.Items),

		Array: &updatesResults.Items,
	}

	err = templates.RenderHTML(w, "updates", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
