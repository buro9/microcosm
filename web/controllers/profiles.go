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

// ProfilesGet will return a page listing profiles
func ProfilesGet(w http.ResponseWriter, req *http.Request) {
	// Query the profiles
	profiles, err := api.GetProfiles(req.Context(), req.URL.Query())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `profiles`,
		Pagination: models.ParsePagination(profiles.Items),

		Profiles: profiles,
	}

	err = templates.RenderHTML(w, "profiles", data)
	if err != nil {
		fmt.Println("could not render profiles")
		w.Write([]byte(err.Error()))
	}
}

// ProfileGet will return a page displaying a single profile
func ProfileGet(w http.ResponseWriter, req *http.Request) {
	// Query the profile
	profileID := asInt64(req, "profileID")
	profile, err := api.GetProfile(req.Context(), profileID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Query the items that they've created
	q := url.Values{}
	q.Add("since", "-1")
	q.Add("type", "conversation")
	q.Add("type", "event")
	q.Add("type", "profile")
	q.Add("type", "huddle")
	q.Add("type", "comment")
	q.Add("authorId", asString(req, "profileID"))
	q.Add("limit", "10")
	q.Add("sort", "date")

	offset := req.URL.Query().Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	searchResults, err := api.DoSearch(req.Context(), q)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Stitch it together for the template
	data := templates.Data{
		Request: req,
		Site:    bag.GetSite(req.Context()),
		User:    bag.GetProfile(req.Context()),
		Section: `profiles`,

		Profile:       profile,
		SearchResults: searchResults,
	}

	err = templates.RenderHTML(w, "profile", data)
	if err != nil {
		fmt.Println("could not render profile")
		w.Write([]byte(err.Error()))
	}
}
