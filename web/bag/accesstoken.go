package bag

import "context"

func SetAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, contextAccessToken, accessToken)
}

func GetAccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(contextAccessToken).(string)
	return accessToken
}
