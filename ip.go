package ui

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

var (
	cloudflareIPv4s = []string{
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"104.16.0.0/12",
		"108.162.192.0/18",
		"131.0.72.0/22",
		"141.101.64.0/18",
		"162.158.0.0/15",
		"172.64.0.0/13",
		"173.245.48.0/20",
		"188.114.96.0/20",
		"190.93.240.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"199.27.128.0/21",
	}

	cloudflareIPv6s = []string{
		"2400:cb00::/32",
		"2405:8100::/32",
		"2405:b500::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2c0f:f248::/32",
		"2a06:98c0::/29",
	}

	cloudflareCIDRv4s    []*net.IPNet
	cloudflareCIDRv6s    []*net.IPNet
	parseCloudFlareCIDRs sync.Once
)

// ipFromRequest returns the IP address for the client that accessed the site
func ipFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("%q is not IP:port", req.RemoteAddr)
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return nil, fmt.Errorf("%q is not IP:port", req.RemoteAddr)
	}

	if isCloudFlareIP(clientIP) {
		cfip := req.Header.Get("CF-Connecting-IP")
		if cfip == "" {
			return nil,
				fmt.Errorf("CF-Connecting-IP not supplied for a CloudFlare IP")
		}

		cfIP := net.ParseIP(cfip)
		if cfIP == nil {
			return nil, fmt.Errorf("%q is not IP:port", cfip)
		}

		return cfIP, nil
	}

	return clientIP, nil
}

// isCloudFlareIP returns true if the given IP address belongs to CloudFlare
func isCloudFlareIP(ip net.IP) bool {
	parseCloudFlareCIDRs.Do(func() {
		for _, s := range cloudflareIPv4s {
			// Ignoring error, hard-coded list are all valid CIDRs and if not
			// this will be nil which we'll check in a moment
			_, cidr, _ := net.ParseCIDR(s)
			if cidr != nil {
				cloudflareCIDRv4s = append(cloudflareCIDRv4s, cidr)
			}
		}
		for _, s := range cloudflareIPv6s {
			// Ignoring error, hard-coded list are all valid CIDRs and if not
			// this will be nil which we'll check in a moment
			_, cidr, _ := net.ParseCIDR(s)
			if cidr != nil {
				cloudflareCIDRv6s = append(cloudflareCIDRv6s, cidr)
			}
		}
	})

	if ip.To4() != nil {
		// IPv6
		for _, network := range cloudflareCIDRv6s {
			if network.Contains(ip) {
				return true
			}
		}
	} else {
		// IPv4
		for _, network := range cloudflareCIDRv4s {
			if network.Contains(ip) {
				return true
			}
		}
	}

	return false
}
