package controllers

import (
	"fmt"
	"net/http"
)

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
