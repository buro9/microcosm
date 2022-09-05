package controllers

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

// MicrocosmGet will fetch the home page
func MicrocosmGet(w http.ResponseWriter, r *http.Request) {
	microcosmID := asInt64(r, "microcosmID")
	microcosm, err := api.GetMicrocosm(r.Context(), microcosmID)
	if err != nil {
		w.Write([]byte(err.Error()))
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
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
