package controllers

import "net/http"

// LogoutPost will remove the session cookie, thus logging the user out
func LogoutPost(w http.ResponseWriter, req *http.Request) {
	var cookie http.Cookie
	cookie.Name = "session"
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Domain = req.Host
	cookie.Path = "/"
	cookie.HttpOnly = true
	if req.URL.Scheme == "https" {
		cookie.Secure = true
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, req, "/", http.StatusFound)
}
