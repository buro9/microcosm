package controllers

import "net/http"

// LogoutPost will remove the session cookie, thus logging the user out
func LogoutPost(w http.ResponseWriter, r *http.Request) {
	var cookie http.Cookie
	cookie.Name = "session"
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Domain = r.Host
	cookie.Path = "/"
	cookie.HttpOnly = true
	if r.URL.Scheme == "https" {
		cookie.Secure = true
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
