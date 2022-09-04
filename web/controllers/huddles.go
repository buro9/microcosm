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
func HuddlesGet(w http.ResponseWriter, req *http.Request) {
	// Set the query options
	q := url.Values{}	
	offset := req.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}
	if req.URL.Query().Get("unread") == "true" {
		q.Add("unread", "true")
	}

	huddles, err := api.GetHuddles(req.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `huddles`,
		Pagination: models.ParsePagination(huddles.Items),

		Huddles: huddles,
	}

	err = templates.RenderHTML(w, "huddles", data)
	if err != nil {
		fmt.Println("could not render huddles")
		w.Write([]byte(err.Error()))
	}
}