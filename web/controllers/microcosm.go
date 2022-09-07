package controllers

import (
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates"
)

// MicrocosmGet will fetch the home page
func MicrocosmGet(w http.ResponseWriter, r *http.Request) {
	microcosmID := asInt64(r, "microcosmID")
	microcosm, status, err := api.GetMicrocosm(r.Context(), microcosmID)
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(microcosm.Items),

		Microcosm: microcosm,
	}

	err = templates.RenderHTML(w, "microcosm", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
