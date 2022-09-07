package errors

import "net/http"

func Render(w http.ResponseWriter, req *http.Request, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Microcosm web client error:\n"))
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}
