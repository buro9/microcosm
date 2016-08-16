package ui

import "net/http"

func homeGet(w http.ResponseWriter, req *http.Request) {
	microcosm, err := microcosms(req.Context())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templateData{
		Request:   req,
		Site:      siteFromContext(req.Context()),
		User:      userFromContext(req.Context()),
		Section:   `home`,
		Microcosm: microcosm,
	}

	err = renderHTMLTemplate(w, "home.tmpl", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
