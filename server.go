package ui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

func ListenAndServe() chan error {
	router := chi.NewRouter()

	// Pages group, handles all routes for pages and defines the appropriate
	// middleware for web pages
	router.Group(func(router chi.Router) {
		router.Use(realIP)
		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.RedirectSlashes)
		router.Use(apiRoot)
		router.Use(forceSSL)
		router.Use(session)

		router.Get("/", homeGet)
	})

	// Static file group, defines minimal middleware
	router.Group(func(router chi.Router) {
		// TODO: Log the static, disabled during dev
		// router.Use(realIP)
		// router.Use(middleware.RequestID)
		// router.Use(middleware.Logger)
		// router.Use(middleware.Recoverer)
		// router.Use(middleware.RedirectSlashes)
		// router.Use(apiRoot)
		router.Use(forceSSL)

		router.Mount("/static", staticFiles())
		// TODO: clear these stubs
		ok := func(w http.ResponseWriter, req *http.Request) { w.Write([]byte("OK")) }

		router.Get("/isogram", ok)
		router.Get("/favicon.ico", ok)
		router.Get("/robots.txt", ok)
	})

	// This is the microcosm client and can work over http as well as https,
	// whilst we'll handle redirecting all *.apidomain.tld to https and likewise
	// for any *Site.ForceSSL to https... we cannot do it for every site as some
	// will be CNAMEd to us and we do not have the certs for their
	// customdomain.tld
	//
	// This means that we serve *everything* over both http and https and we
	// use the forceSSL middleware to use SSL where needed.
	//
	// The by-product of this long-winded explanation is that we listen for both
	// standard http and TLS connections

	// Channel for returning any error out of either of the http or https
	// listeners
	errs := make(chan error)

	go func() {
		log.Printf("Listening for HTTPS on %d ...", *tlsListenPort)
		err := http.ListenAndServeTLS(
			fmt.Sprintf(":%d", *tlsListenPort),
			*certFile,
			*keyFile,
			router,
		)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		log.Printf("Listening for HTTP on %d ...", *listenPort)
		err := http.ListenAndServe(
			fmt.Sprintf(":%d", *listenPort),
			router,
		)
		if err != nil {
			errs <- err
		}
	}()

	return errs
}

func staticFiles() http.Handler {
	router := chi.NewRouter()

	// Do nothing, but implement http.Handler
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, router *http.Request) {
			next.ServeHTTP(w, router)
		})
	})

	// Serve static files
	router.Mount("/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir(*filesPath+"/static/")),
		),
	)

	return router
}
