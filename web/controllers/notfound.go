package controllers

import (
	"fmt"
	"net/http"
)

// NotFound will return a 404 page
func NotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"404 not found for URL %s",
				req.URL.String(),
			),
		),
	)
}

// NotFoundStatic will return a 404 page for a static item
func NotFoundStatic(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"404 not found for URL %s",
				req.URL.String(),
			),
		),
	)
}
