package controllers

import "net/http"

func renderError(w http.ResponseWriter, req *http.Request, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Microcosm web client error:\n"))
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
