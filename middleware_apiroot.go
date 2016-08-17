package ui

import (
	"context"
	"fmt"
	"net/http"
)

// apiRoot is a middleware that populates the context with the root path of the
// API that serves this site. If this cannot be determined then this is not a
// valid Microcosm site and we error out
func apiRoot(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// Get the URL that is the root of the API for this site and store it
		// in the request context
		apiRoot, err := apiRootFromRequest(req)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req = req.WithContext(
			context.WithValue(req.Context(), contextAPIRoot, apiRoot),
		)

		h.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
