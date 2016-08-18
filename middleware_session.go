package ui

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// sessions is a middleware that populates the context with the necessary data
// to complete a session from the perspective of the Microcosm API.
//
// Specifically adds:
//   * Callee IP address to context
//   * Access Token to context (if it is available in the querystring, header or
//     cookie)
//   * Site to context
//   * User to context (if applicable and access token exists and is valid)
//
// This middleware should be inserted last in the middleware stack to ensure
// that information it requires is already available to it (the realIP and the
// apiRoot).
//
// This middleware should *not* be applied to any static files as it does
// perform some processing to fetch the *Site and *User information.
func session(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// Add the IP address to the context as not all funcs that will be
		// called later in the life of this request will be passed the full
		// request
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				contextIP,
				req.RemoteAddr,
			),
		)

		// Get the access_token, if the request has one, and store it in the
		// context
		at := accessTokenFromRequest(req)
		if at != "" {
			req = req.WithContext(
				context.WithValue(req.Context(), contextAccessToken, at),
			)
		}

		// Get the Site based on our knowledge of the API
		site, err := siteFromAPIContext(req.Context())
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req = req.WithContext(
			context.WithValue(req.Context(), contextSite, site),
		)

		// If this site demands SSL and we are not already forcing it, do so
		if site != nil && site.ForceSSL {
			forceSSLHostsLock.RLock()
			_, ok := forceSSLHosts[req.Host]
			forceSSLHostsLock.RUnlock()
			if !ok {
				forceSSLHostsLock.Lock()
				forceSSLHosts[req.Host] = struct{}{}
				forceSSLHostsLock.Unlock()
			}
		}

		// Get the User based on our knowledge of the API
		user, err := userFromAPIContext(req.Context())
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req = req.WithContext(
			context.WithValue(req.Context(), contextUser, user),
		)

		// The IP is stored in the context
		h.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// accessTokenFromRequest returns the access token, if there is one, associated
// with the current request
func accessTokenFromRequest(req *http.Request) string {
	// querystring has precedence
	if at := req.URL.Query().Get("access_token"); at != "" {
		return at
	}

	// then an auth header
	auth := req.Header.Get("Authorisation")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.Replace(auth, "Bearer ", "", 1)
		}
	}

	// finally the cookie
	cookie, _ := req.Cookie("access_token")
	if cookie != nil && cookie.Value != "" && cookie.Secure {
		return cookie.Value
	}

	return ""
}
