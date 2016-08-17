package ui

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	forceSSLHosts     map[string]struct{}
	forceSSLHostsLock sync.RWMutex
)

// forceSSL is a middleware that looks at the request scheme and host to
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
func forceSSL(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if req.TLS == nil {
			if strings.HasSuffix(req.Host, *apiDomain) {
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

	if *tlsListenPort == 443 {
		return fmt.Sprintf(
			"https://%s%s",
			req.Host,
			req.URL.RequestURI(),
		)
	}

	return fmt.Sprintf(
		"https://%s:%d%s",
		req.Host,
		*tlsListenPort,
		req.URL.RequestURI(),
	)
}
