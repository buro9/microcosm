package ui

import "net/http"

func homeGet(w http.ResponseWriter, req *http.Request) {
	ctx, err := newContext(req)

	microcosm, err := microcosms(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templateData{
		Request: req,
		Site:    siteFromContext(ctx),
		Section: `home`,

		Microcosm: microcosm,
	}

	err = renderHTMLTemplate(w, "home.tmpl", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
