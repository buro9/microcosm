package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/opts"
)

// Session is a middleware that populates the context with the necessary data
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
func Session(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// Get the access_token, if the request has one, and store it in the
		// context
		at := accessTokenFromRequest(req)
		if at != "" {
			req = req.WithContext(bag.SetAccessToken(req.Context(), at))
		}

		// Get the Site based on our knowledge of the API
		site, err := api.SiteFromAPIContext(req.Context())
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req = req.WithContext(bag.SetSite(req.Context(), site))

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

		// Get the current profile based on our knowledge of the API
		profile, err := api.ProfileFromAPIContext(req.Context())
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req = req.WithContext(bag.SetProfile(req.Context(), profile))

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
	if cookie, err := req.Cookie("session"); err == nil {
		value := make(map[string]string)
		if err = opts.SecureCookie.Decode("session", cookie.Value, &value); err == nil {
			return value["accessToken"]
		}
	}

	// cookie, _ := req.Cookie("session")
	// if cookie != nil && cookie.Value != "" {
	// 	return cookie.Value
	// }

	return ""
}
