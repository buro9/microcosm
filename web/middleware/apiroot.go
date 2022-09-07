package middleware

import (
	"net/http"

	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
)

// APIRoot is a middleware that populates the context with the root path of the
// API that serves this site. If this cannot be determined then this is not a
// valid Microcosm site and we error out
func APIRoot(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the URL that is the root of the API for this site and store it
		// in the request context
		apiRoot, status, err := api.RootFromRequest(r)
		if err != nil {
			errors.Render(w, r, status, err)
			return
		}

		r = r.WithContext(bag.SetAPIRoot(r.Context(), apiRoot))

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
