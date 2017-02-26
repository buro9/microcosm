package opts

import (
	"flag"
	"os"
	"strings"
	"sync"
)

var parseFlags sync.Once

var (
	// FilesPath is the path to the directory on the file system that contains
	// the static file and template directories
	FilesPath *string

	// Listen is the addr:port that the HTTP server will listen on
	Listen *string

	// TLSListen is the addr:port that the TLS server will listen on
	TLSListen *string

	// CertFile is the path to the certificate that will be used for the TLS
	// connection
	CertFile *string

	// KeyFile is the path to the private key for the CertFile
	KeyFile *string

	// ApiDomain is the top level domain name that serves the api
	ApiDomain *string

	// ClientSecret is the secret that this client uses when talking to the API
	// for exchanging auth credentials
	ClientSecret *string
)

// RegisterFlags adds the flags needed by the UI if they have not already been
// added
func RegisterFlags() {
	parseFlags.Do(
		func() {
			FilesPath = flag.String(
				"files",
				"/srv/microcosm-web",
				"directory that contains the templates and static files",
			)
			if *FilesPath != "" {
				*FilesPath = strings.TrimRight(*FilesPath, "/")
			}

			Listen = flag.String(
				"listen",
				":80",
				"addr:port on which to serve HTTP",
			)

			TLSListen = flag.String(
				"tlsListen",
				":443",
				"addr:port on which to serve HTTPS",
			)

			CertFile = flag.String(
				"certFile",
				"/etc/ssl/certs/microco.sm.crt",
				"path to the TLS certificate file",
			)

			KeyFile = flag.String(
				"keyFile",
				"/etc/ssl/private/microco.sm.key",
				"path to the TLS private key file",
			)

			ApiDomain = flag.String(
				"apiDomain",
				"microco.sm",
				"the .tld that serves the API",
			)

			ClientSecret = flag.String(
				"clientSecret",
				os.Getenv("MICROCOSM_API_CLIENT_SECRET"),
				"the API client secret",
			)
		},
	)
}
