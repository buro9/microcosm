package ui

import (
	"flag"
	"strings"
	"sync"
)

var parseFlags sync.Once

var (
	// filesPath is the path to the directory on the file system that contains
	// the static file and template directories
	filesPath *string

	// listenPort is the port that the server will listen on
	listenPort *int

	// certFile is the path to the certificate that will be used for the TLS
	// connection
	certFile *string

	// keyFile is the path to the private key for the certFile
	keyFile *string

	// apiDomain is the top level domain name that serves the api
	apiDomain *string
)

// RegisterFlags adds the flags needed by the UI if they have not already been
// added
func RegisterFlags() {
	parseFlags.Do(func() {
		filesPath = flag.String(
			"filespath",
			"/srv/microcosm-ui",
			"directory that contains the templates and static files",
		)
		if *filesPath != "" {
			*filesPath = strings.TrimRight(*filesPath, "/")
		}

		listenPort = flag.Int(
			"port",
			443,
			"port on which to serve HTTPS",
		)

		certFile = flag.String(
			"certFile",
			"/etc/ssl/certs/microco.sm.crt",
			"path to the TLS certificate file",
		)

		keyFile = flag.String(
			"keyFile",
			"/etc/ssl/private/microco.sm.key",
			"path to the TLS private key file",
		)

		apiDomain = flag.String(
			"apiDomain",
			"microco.sm",
			"the .tld that serves the API",
		)
	})
}
