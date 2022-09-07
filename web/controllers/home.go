package controllers

import (
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates"
)

// HomeGet will fetch the home page
func HomeGet(w http.ResponseWriter, r *http.Request) {
	rootMicrocosm, status, err := api.GetMicrocosm(r.Context(), 0)
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(rootMicrocosm.Items),

		Microcosm: rootMicrocosm,
	}

	err = templates.RenderHTML(w, "home", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
