package ui

import (
	"fmt"
	"net/http"
)

func registerHandlers() {
	// the static files
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir(*filesPath+"/static/")),
		),
	)

	http.HandleFunc("/", homeGet)

	// TODO: clear these stubs
	ok := func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("OK")) }
	http.HandleFunc("/isogram", ok)
	http.HandleFunc("/favicon.ico", ok)
}

func ListenAndServe() error {
	registerHandlers()

	return http.ListenAndServeTLS(
		fmt.Sprintf(":%d", *listenPort),
		*certFile,
		*keyFile,
		nil,
	)
}
