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

func TodayGet(w http.ResponseWriter, req *http.Request) {
	q := url.Values{}
	q.Add("since", "-1")
	q.Add("type", "conversation")
	q.Add("type", "event")
	q.Add("type", "profile")
	q.Add("type", "huddle")

	offset := req.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	searchResults, err := api.DoSearch(req.Context(), &q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `today`,
		Pagination: models.ParsePagination(searchResults.Items),

		SearchResults: searchResults,
	}

	err = templates.RenderHTML(w, "today", data)
	if err != nil {
		fmt.Println("could not render today")
		w.Write([]byte(err.Error()))
	}
}
