package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates"
)

// HuddlesGet will fetch the home page
func HuddlesGet(w http.ResponseWriter, r *http.Request) {
	user := bag.GetProfile(r.Context())
	if user == nil {
		errors.Render(w, r, http.StatusForbidden, fmt.Errorf(`Need to be signed in to view huddles`))
		return
	}

	// Set the query options
	q := url.Values{}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}
	if r.URL.Query().Get("unread") == "true" {
		q.Add("unread", "true")
	}

	huddles, status, err := api.GetHuddles(r.Context(), q)
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       user,
		Section:    `huddles`,
		Pagination: models.ParsePagination(huddles.Items),

		Huddles: huddles,
	}

	err = templates.RenderHTML(w, "huddles", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
