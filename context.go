package ui

import (
	"context"
	"net"
)

const (
	contextIP = iota
	contextAccessToken
	contextAPIRoot
	contextSite
)

// ipFromContext returns the IP address of the client from the context
func ipFromContext(ctx context.Context) net.IP {
	ip, _ := ctx.Value(contextIP).(net.IP)
	return ip
}

// apiRootFromContext returns the api url for the site associated with the
// current request, i.e. https://subdomain.apidomain.tld/api/v1
func apiRootFromContext(ctx context.Context) string {
	apiRoot, _ := ctx.Value(contextAPIRoot).(string)
	return apiRoot
}

// accessTokenFromContext returns the access token from the context
func accessTokenFromContext(ctx context.Context) string {
	accessToken, _ := ctx.Value(contextAccessToken).(string)
	return accessToken
}

// siteFromContext returns the current site from the context
func siteFromContext(ctx context.Context) *Site {
	site, _ := ctx.Value(contextSite).(*Site)
	return site
}
