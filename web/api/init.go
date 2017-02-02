package api

func init() {
	// Used within apiroot.go
	cnameToAPIRootLock.Lock()
	cnameToAPIRoot = make(map[string]string)
	cnameToAPIRootLock.Unlock()
}
