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
	// Knowing this, we're just going to pass this straight through whereas
	// everywhere else we effectively sanitise the input by checking every arg
	// as that would improve cacheability of the underlying resource, but here
	// we're allowing something deep inside the API to do that.
	q := r.URL.Query()

	if q.Has("defaults") {
		q.Set("inTitle", "true")
		q.Set("sort", "date")
	}

	searchResults, status, err := api.DoSearch(r.Context(), q)
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `search`,
		Query:      &searchResults.Query,
		Pagination: models.ParsePagination(searchResults.Items),

		SearchResults: searchResults,
	}

	err = templates.RenderHTML(w, "search", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
