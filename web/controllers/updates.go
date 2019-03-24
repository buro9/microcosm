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
func UpdatesGet(w http.ResponseWriter, req *http.Request) {
	q := url.Values{}
	offset := req.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	updatesResults, err := api.GetUpdates(req.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `updates`,
		Pagination: models.ParsePagination(updatesResults.Items),

		Array: &updatesResults.Items,
	}

	err = templates.RenderHTML(w, "updates", data)
	if err != nil {
		fmt.Println("could not render updates")
		w.Write([]byte(err.Error()))
	}
}
