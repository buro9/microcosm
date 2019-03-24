package controllers

import (
	"net/http"
	"time"

	"github.com/buro9/microcosm/web/api"
)

// Auth0LoginGet will attempt to log the user in using Auth0 and then set
// the session cookie for the current user before redirecting to the destination
func Auth0LoginGet(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	targetURL := q.Get("target_url")

	accessToken, err := api.Auth0Login(req.Context(), code, state)
	if err != nil {
		renderError(w, req, err)
		return
	}

	if accessToken == "" {
		renderError(w, req, api.ErrAccessTokenExpected)
		return
	}

	if targetURL == "" {
		if state != "" {
			targetURL = state
		} else {
			targetURL = `/`
		}
	}

	var cookie http.Cookie
	cookie.Name = "session"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add(time.Hour * 24 * 365)
	cookie.Domain = req.Host
	cookie.Path = "/"
	cookie.HttpOnly = true
	if req.URL.Scheme == "https" {
		cookie.Secure = true
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, req, targetURL, http.StatusFound)
}
