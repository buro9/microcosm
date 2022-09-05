package controllers

import (
	"fmt"
	"net/http"
)

// NotFound will return a 404 page
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"404 not found for URL %s",
				r.URL.String(),
			),
		),
	)
}

// NotFoundStatic will return a 404 page for a static item
func NotFoundStatic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"404 not found for URL %s",
				r.URL.String(),
			),
		),
	)
}
