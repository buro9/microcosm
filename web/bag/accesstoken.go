package bag

import "context"

// SetAccessToken puts the access token within the context
func SetAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, contextAccessToken, accessToken)
}

// GetAccessToken fetches the access token from the context
func GetAccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(contextAccessToken).(string)
	return accessToken
}
