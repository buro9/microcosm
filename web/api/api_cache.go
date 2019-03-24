package api

import (
	"github.com/gregjones/httpcache/memcache"
)

var apiCache *memcache.Cache

// NewCache will create an API cache for a given memcached server address
func NewCache(addr string) {
	apiCache = memcache.New(addr)
}
