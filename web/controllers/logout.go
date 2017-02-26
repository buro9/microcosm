package controllers

import "net/http"

func LogoutPost(w http.ResponseWriter, req *http.Request) {
	var cookie http.Cookie
	cookie.Name = "access_token"
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
