package ui

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	contextIP = iota
	contextAccessToken
	contextAPIRoot
	contextSite
)

// newContext constructs a context with all of the applicable knowledge needed
// to produce a page and perform any other actions
func newContext(req *http.Request) (context.Context, error) {
	ctx := context.Background()

	// Store the IP address of the user
	ip, err := ipFromRequest(req)
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, contextIP, ip)

	// Get the access_token, if the request has one
	at := accessTokenFromRequest(req)
	if at != "" {
		ctx = context.WithValue(ctx, contextAccessToken, ip)
	}

	// Get the URL that is the root of the API for this site
	apiRoot, err := apiRootFromRequest(req)
	if err != nil {
		fmt.Println(err.Error())
		return ctx, err
	}
	ctx = context.WithValue(ctx, contextAPIRoot, apiRoot)

	// Get the Site based on our knowledge of the API
	site, err := siteFromAPIContext(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return ctx, err
	}
	ctx = context.WithValue(ctx, contextSite, site)

	return ctx, nil
}

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

// accessTokenFromRequest returns the access token, if there is one, associated
// with the current request
func accessTokenFromRequest(req *http.Request) string {
	// querystring has precedence
	if at := req.URL.Query().Get("access_token"); at != "" {
		return at
	}

	// then an auth header
	auth := req.Header.Get("Authorisation")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.Replace(auth, "Bearer ", "", 1)
		}
	}

	// finally the cookie
	cookie, _ := req.Cookie("access_token")
	if cookie != nil && cookie.Value != "" && cookie.Secure {
		return cookie.Value
	}

	return ""
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
