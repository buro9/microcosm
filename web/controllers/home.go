package controllers

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

// HomeGet will fetch the home page
func HomeGet(w http.ResponseWriter, r *http.Request) {
	rootMicrocosm, err := api.GetMicrocosm(r.Context(), 0)
	if err != nil {
		w.Write([]byte(err.Error()))
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
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
