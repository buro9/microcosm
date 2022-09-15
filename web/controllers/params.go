package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func asInt64(r *http.Request, name string) int64 {
	i, _ := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	return i
}

func asString(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
