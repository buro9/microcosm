package controllers

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

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
