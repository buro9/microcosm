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

// HuddlesGet will fetch the home page
func HuddlesGet(w http.ResponseWriter, r *http.Request) {
	// Set the query options
	q := url.Values{}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}
	if r.URL.Query().Get("unread") == "true" {
		q.Add("unread", "true")
	}

	huddles, err := api.GetHuddles(r.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `huddles`,
		Pagination: models.ParsePagination(huddles.Items),

		Huddles: huddles,
	}

	err = templates.RenderHTML(w, "huddles", data)
	if err != nil {
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
