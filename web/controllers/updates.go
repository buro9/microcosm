package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

// UpdatesGet will return the updates page
func UpdatesGet(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	updatesResults, err := api.GetUpdates(r.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
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
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
