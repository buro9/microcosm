package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/buro9/microcosm/web/bag"
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

	cfConnectingIP = http.CanonicalHeaderKey("CF-Connecting-IP")
	xRealIP        = http.CanonicalHeaderKey("X-Real-IP")
)

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the CF-Connecting-IP header or the X-Real-IP header
// (in that order).
//
// This middleware should be inserted as the first in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if rip := ipFromRequest(req); rip != nil {
			req.RemoteAddr = rip.String()
		}

		// The IP is stored in the context as not all funcs that will be
		// called later in the life of this request will be passed the full
		// request
		req = req.WithContext(bag.SetIP(req.Context(), req.RemoteAddr))

		h.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// ipFromRequest returns the IP address for the client that accessed the site
func ipFromRequest(req *http.Request) net.IP {
	var realIP net.IP

	// X-Real-IP if supplied
	if xrip := req.Header.Get(xRealIP); xrip != "" {
		if clientIP := net.ParseIP(xrip); clientIP != nil {
			realIP = clientIP
		}
	}

	if realIP == nil {
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		raip := net.ParseIP(ip)
		if raip == nil {
			// Neither X-Real-IP or req.RemoteAddr were parsable, return nothing
			// and leave everything unchanged, this is weird
			return nil
		}
		realIP = raip
	}

	if isCloudFlareIP(realIP) {
		// We only trust this header when we're behind a CloudFlare IP, and we
		// are... so the CF-Connecting-IP actually holds the realIP
		cfip := req.Header.Get(cfConnectingIP)
		if cfip == "" {
			fmt.Printf(
				"CF-Connecting-IP not supplied for CloudFlare IP %s\n",
				realIP.String(),
			)
			return nil
		}

		cfIP := net.ParseIP(cfip)
		if cfIP == nil {
			fmt.Printf(
				"CF-Connecting-IP supplied an invalid IP %s\n",
				cfip,
			)
			return nil
		}

		realIP = cfIP
	}

	return realIP
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
