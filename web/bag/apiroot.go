package bag

import "context"

func SetAPIRoot(ctx context.Context, apiRoot string) context.Context {
	return context.WithValue(ctx, contextAPIRoot, apiRoot)
}

// GetAPIRoot returns the api url for the site associated with the
// current request, i.e. https://subdomain.apidomain.tld/api/v1
func GetAPIRoot(ctx context.Context) string {
	apiRoot, _ := ctx.Value(contextAPIRoot).(string)
	return apiRoot
}
