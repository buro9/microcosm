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
func ProfilesGet(w http.ResponseWriter, r *http.Request) {
	// Query the profiles
	profiles, err := api.GetProfiles(r.Context(), r.URL.Query())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `profiles`,
		Pagination: models.ParsePagination(profiles.Items),

		Profiles: profiles,
	}

	err = templates.RenderHTML(w, "profiles", data)
	if err != nil {
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}

// ProfileGet will return a page displaying a single profile
func ProfileGet(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	// Query the profile
	var (
		profile    *models.Profile
		profileErr error
	)

	profileID := asInt64(r, "profileID")

	wg.Add(1)
	go func(ctx context.Context, profileID int64) {
		defer wg.Done()
		profile, profileErr = api.GetProfile(r.Context(), profileID)
	}(r.Context(), profileID)

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
	q.Add("authorId", asString(r, "profileID"))
	q.Add("limit", "10")
	q.Add("sort", "date")

	wg.Add(1)
	go func(ctx context.Context, q url.Values) {
		defer wg.Done()
		searchResults, searchErr = api.DoSearch(r.Context(), q)
	}(r.Context(), q)

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
		Request: r,
		Site:    bag.GetSite(r.Context()),
		User:    bag.GetProfile(r.Context()),
		Section: `profiles`,

		Profile:       profile,
		SearchResults: searchResults,
	}

	err := templates.RenderHTML(w, "profile", data)
	if err != nil {
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
