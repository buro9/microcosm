package bag

import "context"

// SetAPIRoot puts the API URI for the site associated with the current request
// into the context
func SetAPIRoot(ctx context.Context, apiRoot string) context.Context {
	return context.WithValue(ctx, contextAPIRoot, apiRoot)
}

// GetAPIRoot fetches the api url for the site associated with the
// current request, i.e. https://subdomain.apidomain.tld/api/v1
func GetAPIRoot(ctx context.Context) string {
	// apiRoot, _ := ctx.Value(contextAPIRoot).(string)
	// return apiRoot

	// TODO: Uncomment the above
	return "https://gfora.microco.sm/api/v1"
}
