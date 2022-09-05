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

// TodayGet will return the today page
func TodayGet(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Add("since", "-1")
	q.Add("type", "conversation")
	q.Add("type", "event")
	q.Add("type", "profile")
	q.Add("type", "huddle")

	offset := r.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	searchResults, err := api.DoSearch(r.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `today`,
		Pagination: models.ParsePagination(searchResults.Items),

		SearchResults: searchResults,
	}

	err = templates.RenderHTML(w, "today", data)
	if err != nil {
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
