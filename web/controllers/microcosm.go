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
func MicrocosmGet(w http.ResponseWriter, req *http.Request) {
	microcosmID := asInt64(req, "microcosmID")
	microcosm, err := api.GetMicrocosm(req.Context(), microcosmID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(microcosm.Items),

		Microcosm: microcosm,
	}

	err = templates.RenderHTML(w, "home", data)
	if err != nil {
		fmt.Println("could not render home")
		w.Write([]byte(err.Error()))
	}
}
