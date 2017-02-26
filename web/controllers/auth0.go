package controllers

import (
	"net/http"
	"time"

	"github.com/buro9/microcosm/web/api"
)

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
	cookie.Name = "access_token"
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
