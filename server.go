package ui

import (
	"fmt"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func ListenAndServe() error {
	r := chi.NewRouter()

	// Pages group, handles all routes for pages and defines the appropriate
	// middleware for web pages
	r.Group(func(r chi.Router) {
		r.Use(realIP)
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RedirectSlashes)
		r.Use(session)

		r.Get("/", homeGet)
	})

	// Static file group, defines minimal middleware
	r.Group(func(r chi.Router) {
		// TODO: Log the static, disabled during dev
		// r.Use(realIP)
		// r.Use(middleware.RequestID)
		// r.Use(middleware.Logger)
		// r.Use(middleware.Recoverer)
		// r.Use(middleware.RedirectSlashes)

		r.Mount("/static", staticFiles())
		// TODO: clear these stubs
		ok := func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("OK")) }

		r.Get("/isogram", ok)
		r.Get("/favicon.ico", ok)
		r.Get("/robots.txt", ok)
	})

	return http.ListenAndServeTLS(
		fmt.Sprintf(":%d", *listenPort),
		*certFile,
		*keyFile,
		r,
	)
}

func staticFiles() http.Handler {
	r := chi.NewRouter()

	// Do nothing, but implement http.Handler
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})

	// Serve static files
	r.Mount("/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir(*filesPath+"/static/")),
		),
	)

	return r
}
