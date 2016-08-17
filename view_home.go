package ui

import (
	"fmt"
	"net/http"
)

func homeGet(w http.ResponseWriter, req *http.Request) {
	microcosm, err := microcosms(req.Context())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templateData{
		Request:    req,
		Site:       siteFromContext(req.Context()),
		User:       userFromContext(req.Context()),
		Section:    `home`,
		Pagination: parsePagination(microcosm.Items),

		Microcosm: microcosm,
	}

	err = renderHTMLTemplate(w, "home.tmpl", data)
	if err != nil {
		fmt.Println("could not render home.tmpl")
		w.Write([]byte(err.Error()))
	}
}
