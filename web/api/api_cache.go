package api

import (
	"github.com/gregjones/httpcache/memcache"
)

var apiCache *memcache.Cache

func NewAPICache(addr string) {
	apiCache = memcache.New(addr)
}
