package ui

import "github.com/oxtoacart/bpool"

func init() {
	// Used within templates.go
	bufpool = bpool.NewBufferPool(64)

	// Used within context.go
	cnameToAPIRootLock.Lock()
	cnameToAPIRoot = make(map[string]string)
	cnameToAPIRootLock.Unlock()

	// Used within middleware_forcessl.go
	forceSSLHostsLock.Lock()
	forceSSLHosts = make(map[string]struct{})
	forceSSLHostsLock.Unlock()
}
