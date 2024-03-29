package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gorilla/csrf"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/buro9/microcosm/web/controllers"
	mm "github.com/buro9/microcosm/web/middleware"
	"github.com/buro9/microcosm/web/opts"
)

// ListenAndServe will run the web server
func ListenAndServe(version string) chan error {
	router := chi.NewRouter()

	csrfKeyString := *opts.CsrfKey
	csrfKeyBytes := []byte(csrfKeyString)

	csrfMiddleware :=
		csrf.Protect(csrfKeyBytes,
			csrf.Secure(false),
			csrf.FieldName("csrfmiddlewaretoken"))

	// Pages group, handles all routes for pages and defines the appropriate
	// middleware for web pages
	router.Group(func(router chi.Router) {
		router.Use(mm.RealIP)
		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(mm.APIRoot)
		router.Use(mm.ForceSSL)
		router.Use(mm.Session)
		router.Use(csrfMiddleware)

		// Login
		router.Get(`/`, controllers.HomeGet)
		router.Get(`/auth0login/`, controllers.Auth0LoginGet)
		router.Post(`/logout/`, controllers.LogoutPost)

		// Comments
		//
		// Note that these can be added to anything; Conversations, Events, Profile, etc.
		router.Post(`/comments/create`, controllers.CommentCreate)

		// Conversations, a.k.a Threads
		router.Get(`/conversations/{conversationID:[1-9][0-9]+}/`, controllers.ConversationGet)
		router.Get(`/conversations/{conversationID:[1-9][0-9]+}/{jumpTo:newest}/`, controllers.ConversationGet)

		// Huddles, a.k.a Private Messages or Direct Messages
		router.Get(`/huddles/`, controllers.HuddlesGet)

		// Microcosms, a.k.a Forums
		router.Get(`/microcosms/{microcosmID:[1-9][0-9]+}/`, controllers.MicrocosmGet)

		// Profiles
		router.Get(`/profiles/{profileID:[1-9][0-9]+}/`, controllers.ProfileGet)
		router.Get(`/profiles/`, controllers.ProfilesGet)

		// Search
		router.Get(`/search/`, controllers.SearchGet)

		// Today
		router.Get(`/today/`, controllers.TodayGet)

		// Updates, a.k.a Watched items, Followed items
		router.Get(`/updates/`, controllers.UpdatesGet)

		router.NotFound(controllers.NotFound)
	})

	// Static file group, defines minimal middleware
	router.Group(func(router chi.Router) {
		router.Use(mm.RealIP)
		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(mm.ForceSSL)

		router.Mount(`/static`, staticFiles(version))

		router.Get(`/favicon.ico`, func(w http.ResponseWriter, r *http.Request) {
			file, err := inlinedFiles.ReadFile("static/favicon.ico")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(
					[]byte(fmt.Sprintf("500 server error: %s", err.Error())),
				)
				return
			}

			if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && lastModified.Before(t.Add(1*time.Second)) {
				w.WriteHeader(304)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set(`Content-Type`, `image/png`)
			w.Header().Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))
			w.Header().Set("Cache-Control", "no-cache")
			w.Write(file)

			return
		})

		router.Get(`/robots.txt`, func(w http.ResponseWriter, r *http.Request) {
			file, err := inlinedFiles.ReadFile("static/robots.txt")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(
					[]byte(fmt.Sprintf("500 server error: %s", err.Error())),
				)
				return
			}

			if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && lastModified.Before(t.Add(1*time.Second)) {
				w.WriteHeader(304)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set(`Content-Type`, `text/plain`)
			w.Header().Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))
			w.Header().Set("Cache-Control", "no-cache")
			w.Write(file)
			return
		})

		router.NotFound(controllers.NotFoundStatic)
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
		log.Printf("Listening for HTTPS on %s ...", *opts.TLSListen)
		err := http.ListenAndServeTLS(
			*opts.TLSListen,
			*opts.CertFile,
			*opts.KeyFile,
			router,
		)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		log.Printf("Listening for HTTP on %s ...", *opts.Listen)
		err := http.ListenAndServe(
			*opts.Listen,
			router,
		)
		if err != nil {
			errs <- err
		}
	}()

	return errs
}

//go:embed static/*
var inlinedFiles embed.FS
var lastModified time.Time = time.Now()

func staticFiles(version string) http.Handler {
	router := chi.NewRouter()

	// Do nothing, but implement http.Handler
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, version+"/", "", 1)

			switch {
			case strings.HasSuffix(r.URL.Path, `.apng`):
				w.Header().Set(`Content-Type`, `image/apng`)

			case strings.HasSuffix(r.URL.Path, `.avif`):
				w.Header().Set(`Content-Type`, `image/avif`)

			case strings.HasSuffix(r.URL.Path, `.css`):
				w.Header().Set(`Content-Type`, `text/css`)

			case strings.HasSuffix(r.URL.Path, `.gif`):
				w.Header().Set(`Content-Type`, `image/gif`)

			case strings.HasSuffix(r.URL.Path, `.jpg`):
				w.Header().Set(`Content-Type`, `image/jpeg`)

			case strings.HasSuffix(r.URL.Path, `.jpeg`):
				w.Header().Set(`Content-Type`, `image/jpeg`)

			case strings.HasSuffix(r.URL.Path, `.jpe`):
				w.Header().Set(`Content-Type`, `image/jpeg`)

			case strings.HasSuffix(r.URL.Path, `.jif`):
				w.Header().Set(`Content-Type`, `image/jpeg`)

			case strings.HasSuffix(r.URL.Path, `.jfif`):
				w.Header().Set(`Content-Type`, `image/jpeg`)

			case strings.HasSuffix(r.URL.Path, `.js`):
				w.Header().Set(`Content-Type`, `text/javascript`)

			case strings.HasSuffix(r.URL.Path, `.png`):
				w.Header().Set(`Content-Type`, `image/png`)

			case strings.HasSuffix(r.URL.Path, `.svg`):
				w.Header().Set(`Content-Type`, `image/svg+xml`)

			case strings.HasSuffix(r.URL.Path, `.webp`):
				w.Header().Set(`Content-Type`, `image/webp`)
			}

			if *opts.IsDevelopment == false {
				w.Header().Set("Cache-Control", "public,max-age=604800")
			} else {
				w.Header().Set("Cache-Control", "no-cache")
			}

			next.ServeHTTP(w, r)
		})
	})

	if *opts.IsDevelopment {
		// Serve static files from disk
		cwd, _ := os.Getwd()
		fileSystem := http.Dir(path.Join(cwd, "web/server/static"))
		router.Mount(`/`,
			http.StripPrefix("/static/", http.FileServer(fileSystem)),
		)
	} else {
		// Serve static files from embedded artefacts
		router.Mount(`/`,
			http.FileServer(http.FS(inlinedFiles)),
		)
	}

	return router
}
