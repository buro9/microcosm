package controllers

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

func asInt64(req *http.Request, name string) int64 {
	i, _ := strconv.ParseInt(chi.URLParam(req, name), 10, 64)
	return i
}

func asString(req *http.Request, name string) string {
	return chi.URLParam(req, name)
}
