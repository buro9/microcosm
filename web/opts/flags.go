package opts

import (
	"flag"
	"strings"
	"sync"
)

var parseFlags sync.Once

var (
	// FilesPath is the path to the directory on the file system that contains
	// the static file and template directories
	FilesPath *string

	// ListenPort is the port that the HTTP server will listen on
	ListenPort *int

	// TLSListenPort is the port that the TLS server will listen on
	TLSListenPort *int

	// CertFile is the path to the certificate that will be used for the TLS
	// connection
	CertFile *string

	// KeyFile is the path to the private key for the CertFile
	KeyFile *string

	// ApiDomain is the top level domain name that serves the api
	ApiDomain *string
)

// RegisterFlags adds the flags needed by the UI if they have not already been
// added
func RegisterFlags() {
	parseFlags.Do(
		func() {
			FilesPath = flag.String(
				"Filespath",
				"/srv/microcosm-web",
				"directory that contains the templates and static files",
			)
			if *FilesPath != "" {
				*FilesPath = strings.TrimRight(*FilesPath, "/")
			}

			ListenPort = flag.Int(
				"port",
				80,
				"port on which to serve HTTP",
			)

			TLSListenPort = flag.Int(
				"tlsPort",
				443,
				"port on which to serve HTTPS",
			)

			CertFile = flag.String(
				"CertFile",
				"/etc/ssl/certs/microco.sm.crt",
				"path to the TLS certificate file",
			)

			KeyFile = flag.String(
				"KeyFile",
				"/etc/ssl/private/microco.sm.key",
				"path to the TLS private key file",
			)

			ApiDomain = flag.String(
				"apiDomain",
				"microco.sm",
				"the .tld that serves the API",
			)
		},
	)
}
