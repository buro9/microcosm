package middleware

func init() {
	// Used within forcessl.go
	forceSSLHostsLock.Lock()
	forceSSLHosts = make(map[string]struct{})
	forceSSLHostsLock.Unlock()
}
