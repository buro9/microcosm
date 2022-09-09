package controllers

import (
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates"
)

// SearchGet will return search results
func SearchGet(w http.ResponseWriter, r *http.Request) {
	// We're not sanitising the input here as the API does that very thoroughly:
	// https://github.com/microcosm-cc/microcosm/blob/main/models/search_query.go
	//
	// Knowing this, we're just going to pass this straight through.
	searchResults, status, err := api.DoSearch(r.Context(), r.URL.Query())
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `search`,
		Pagination: models.ParsePagination(searchResults.Items),

		SearchResults: searchResults,
	}

	err = templates.RenderHTML(w, "search", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
