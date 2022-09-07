package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/buro9/microcosm/web/opts"
)

var (
	forceSSLHosts     map[string]struct{}
	forceSSLHostsLock sync.RWMutex
)

// ForceSSL is a middleware that looks at the request scheme and host to
// determine whether this is over http and should be redirected over https.
//
// The rules for this are:
// 1. If r.URL.Scheme == https do nothing.
// 2. If r.URL.Host == *.apidomain, redirect to https.
// 3. If r.URL.Host exists in forceSSLHosts, redirect to https.
//
// forceSSLHosts is loaded by virtue of the session middleware fetching
// knowledge of the site and then populating the forceSSLHosts if
// Site.ForceSSL is true
func ForceSSL(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			if strings.HasSuffix(r.Host, *opts.APIDomain) {
				http.Redirect(
					w,
					r,
					redirectURLtoTLS(r),
					http.StatusMovedPermanently,
				)
				return
			}

			forceSSLHostsLock.RLock()
			_, ok := forceSSLHosts[r.Host]
			forceSSLHostsLock.RUnlock()

			if ok {
				http.Redirect(
					w,
					r,
					redirectURLtoTLS(r),
					http.StatusMovedPermanently,
				)
				return
			}
		}

		// SSL not being forced, serve the content
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func redirectURLtoTLS(r *http.Request) string {
	if r.TLS != nil {
		return r.URL.String()
	}

	if strings.Contains(*opts.TLSListen, ":443") {
		return fmt.Sprintf(
			"https://%s%s",
			r.Host,
			r.URL.RequestURI(),
		)
	}

	addrPort := strings.Split(*opts.TLSListen, ":")
	if len(addrPort) != 2 || addrPort[1] == "443" {
		return fmt.Sprintf(
			"https://%s%s",
			r.Host,
			r.URL.RequestURI(),
		)
	}

	return fmt.Sprintf(
		"https://%s:%s%s",
		r.Host,
		addrPort[1],
		r.URL.RequestURI(),
	)
}
