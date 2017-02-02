package bag

import "context"

func SetIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, contextAPIRoot, ip)
}

func GetIP(ctx context.Context) string {
	ip, _ := ctx.Value(contextIP).(string)
	return ip
}
