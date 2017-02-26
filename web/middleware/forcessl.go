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
// 1. If req.URL.Scheme == https do nothing.
// 2. If req.URL.Host == *.apidomain, redirect to https.
// 3. If req.URL.Host exists in forceSSLHosts, redirect to https.
//
// forceSSLHosts is loaded by virtue of the session middleware fetching
// knowledge of the site and then populating the forceSSLHosts if
// Site.ForceSSL is true
func ForceSSL(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.TLS == nil {
			if strings.HasSuffix(req.Host, *opts.ApiDomain) {
				http.Redirect(
					w,
					req,
					redirectURLtoTLS(req),
					http.StatusMovedPermanently,
				)
				return
			}

			forceSSLHostsLock.RLock()
			_, ok := forceSSLHosts[req.Host]
			forceSSLHostsLock.RUnlock()

			if ok {
				http.Redirect(
					w,
					req,
					redirectURLtoTLS(req),
					http.StatusMovedPermanently,
				)
				return
			}
		}

		// SSL not being forced, serve the content
		h.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

func redirectURLtoTLS(req *http.Request) string {
	if req.TLS != nil {
		return req.URL.String()
	}

	if strings.Contains(*opts.TLSListen, ":443") {
		return fmt.Sprintf(
			"https://%s%s",
			req.Host,
			req.URL.RequestURI(),
		)
	}

	addrPort := strings.Split(*opts.TLSListen, ":")
	if len(addrPort) != 2 || addrPort[1] == "443" {
		return fmt.Sprintf(
			"https://%s%s",
			req.Host,
			req.URL.RequestURI(),
		)
	}

	return fmt.Sprintf(
		"https://%s:%d%s",
		req.Host,
		addrPort[1],
		req.URL.RequestURI(),
	)
}
