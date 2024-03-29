package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/opts"
)

// Auth0LoginGet will attempt to log the user in using Auth0 and then set
// the session cookie for the current user before redirecting to the destination
func Auth0LoginGet(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	code := q.Get("code")
	state := q.Get("state")
	targetURL := q.Get("target_url")

	accessToken, status, err := api.Auth0Login(req.Context(), code, state)
	if err != nil {
		errors.Render(w, req, status, err)
		return
	}

	if accessToken == "" {
		errors.Render(w, req, status, errors.ErrAccessTokenExpected)
		return
	}

	if targetURL == "" {
		if state != "" {
			targetURL = state
		} else {
			targetURL = `/`
		}
	}

	value := map[string]string{
		"accessToken": accessToken,
	}
	if opts.SecureCookie == nil {
		errors.Render(w, req, status, fmt.Errorf("SecureCookie must exist"))
		return
	}
	encoded, err := (*opts.SecureCookie).Encode("session", value)
	if err != nil {
		errors.Render(w, req, status, err)
		return
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		Domain:   req.Host,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, req, targetURL, http.StatusFound)
}
