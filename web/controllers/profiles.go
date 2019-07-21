package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

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
	var wg sync.WaitGroup

	// Query the profile
	var (
		profile    *models.Profile
		profileErr error
	)

	profileID := asInt64(req, "profileID")

	wg.Add(1)
	go func(ctx context.Context, profileID int64) {
		defer wg.Done()
		profile, profileErr = api.GetProfile(req.Context(), profileID)
	}(req.Context(), profileID)

	// Query the items that they've created
	var (
		searchResults *models.SearchResults
		searchErr     error
	)

	q := url.Values{}
	q.Add("type", "conversation")
	q.Add("type", "event")
	q.Add("type", "profile")
	q.Add("type", "huddle")
	q.Add("type", "comment")
	q.Add("authorId", asString(req, "profileID"))
	q.Add("limit", "10")
	q.Add("sort", "date")

	wg.Add(1)
	go func(ctx context.Context, q url.Values) {
		defer wg.Done()
		searchResults, searchErr = api.DoSearch(req.Context(), q)
	}(req.Context(), q)

	// Wait for all queries and check for errors
	wg.Wait()

	if profileErr != nil {
		w.Write([]byte(profileErr.Error()))
		return
	}

	if searchErr != nil {
		w.Write([]byte(searchErr.Error()))
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

	err := templates.RenderHTML(w, "profile", data)
	if err != nil {
		fmt.Println("could not render profile")
		w.Write([]byte(err.Error()))
	}
}
