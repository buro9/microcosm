package bag

import "context"

// SetIP puts the current client IP into the context
func SetIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, contextAPIRoot, ip)
}

// GetIP fetches the current client IP from the context
func GetIP(ctx context.Context) string {
	ip, _ := ctx.Value(contextIP).(string)
	return ip
}
