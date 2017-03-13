package opts

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	envPrefix = "MICROCOSM_WEB_"

	defaultFilesPath    = "/srv/microcosm-web"
	defaultListen       = ":80"
	defaultTLSListen    = ":443"
	defaultCertFile     = "/etc/ssl/certs/microco.sm.crt"
	defaultKeyFile      = "/etc/ssl/private/microco.sm.key"
	defaultAPIDomain    = "microco.sm"
	defaultClientSecret = ""
	defaultMemcacheAddr = "localhost:11211"
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

	// APIDomain is the top level domain name that serves the api
	APIDomain *string

	// ClientSecret is the secret that this client uses when talking to the API
	// for exchanging auth credentials
	ClientSecret *string

	// MemcacheAddr contains the connection information for memcached, typically
	// the address localhost:11211
	MemcacheAddr *string
)

// RegisterFlags adds the flags needed by the UI if they have not already been
// added.
//
// Note: Defaults are empty as we will try os.Getenv for the values if they are
// provided by the command line, before finally setting defaults internally.
func RegisterFlags() {
	parseFlags.Do(
		func() {
			FilesPath = flag.String(
				"microcosmWebFiles",
				"",
				`directory that contains the templates and static files
	alternatively $`+envPrefix+`FILES
	(default "`+defaultFilesPath+`")`,
			)
			if *FilesPath != "" {
				*FilesPath = strings.TrimRight(*FilesPath, "/")
			}

			Listen = flag.String(
				"microcosmWebListen",
				"",
				`addr:port on which to serve HTTP
	alternatively $`+envPrefix+`LISTEN
	(default "`+defaultListen+`")`,
			)

			TLSListen = flag.String(
				"microcosmWebTLSListen",
				"",
				`addr:port on which to serve HTTPS
	alternatively $`+envPrefix+`TLS_LISTEN
	(default "`+defaultTLSListen+`")`,
			)

			CertFile = flag.String(
				"microcosmWebCertFile",
				"",
				`path to the TLS certificate file
	alternatively $`+envPrefix+`CERT_FILE
	(default "`+defaultCertFile+`")`,
			)

			KeyFile = flag.String(
				"microcosmWebKeyFile",
				"",
				`path to the TLS private key file
	alternatively $`+envPrefix+`KEY_FILE
	(default "`+defaultKeyFile+`")`,
			)

			APIDomain = flag.String(
				"microcosmWebAPIDomain",
				"",
				`the .tld that serves the API
	alternatively $`+envPrefix+`API_DOMAIN
	(default "`+defaultAPIDomain+`")`,
			)

			ClientSecret = flag.String(
				"microcosmWebAPIClientSecret",
				"",
				`the API client secret
	alternatively $`+envPrefix+`API_CLIENT_SECRET
	(default "`+defaultClientSecret+`")`,
			)

			MemcacheAddr = flag.String(
				"microcosmWebMemcacheAddr",
				"",
				`the API client secret
	alternatively $`+envPrefix+`MEMCACHE_ADDR
	(default "`+defaultClientSecret+`")`,
			)
		},
	)
}

// ValidateFlags will check every flag and if the command line args provided
// no values, then the environment is checked and if those provided no values
// then the defaults are used. If this is still insufficient an error is
// returned.
//
// Some attempt is also made to ensure that the values are correct, i.e. that
// paths are readable.
func ValidateFlags() error {
	// FilesPath
	if FilesPath == nil || *FilesPath == "" {
		path := os.Getenv(envPrefix + "FILES")
		FilesPath = &path
	}
	if *FilesPath == "" {
		path := defaultFilesPath
		FilesPath = &path
	}
	if strings.HasSuffix(*FilesPath, `/`) {
		*FilesPath = strings.TrimSuffix(*FilesPath, `/`)
	}
	if _, err := os.Stat(*FilesPath); err != nil {
		fmt.Println(err.Error())
		return ErrFilesPathNotReadable
	}

	// Listen
	if Listen == nil || *Listen == "" {
		listen := os.Getenv(envPrefix + "LISTEN")
		Listen = &listen
	}
	if *Listen == "" {
		listen := defaultListen
		Listen = &listen
	}

	// TLSListen
	if TLSListen == nil || *TLSListen == "" {
		tlsListen := os.Getenv(envPrefix + "TLS_LISTEN")
		TLSListen = &tlsListen
	}
	if *TLSListen == "" {
		tlsListen := defaultTLSListen
		TLSListen = &tlsListen
	}

	// CertFile
	if CertFile == nil || *CertFile == "" {
		certFile := os.Getenv(envPrefix + "CERT_FILE")
		CertFile = &certFile
	}
	if *CertFile == "" {
		certFile := defaultCertFile
		CertFile = &certFile
	}

	// KeyFile
	if KeyFile == nil || *KeyFile == "" {
		keyFile := os.Getenv(envPrefix + "KEY_FILE")
		KeyFile = &keyFile
	}
	if *KeyFile == "" {
		keyFile := defaultKeyFile
		KeyFile = &keyFile
	}

	// APIDomain
	if APIDomain == nil || *APIDomain == "" {
		apiDomain := os.Getenv(envPrefix + "API_DOMAIN")
		APIDomain = &apiDomain
	}
	if *APIDomain == "" {
		apiDomain := defaultAPIDomain
		APIDomain = &apiDomain
	}

	// ClientSecret
	if ClientSecret == nil || *ClientSecret == "" {
		clientSecret := os.Getenv(envPrefix + "API_CLIENT_SECRET")
		ClientSecret = &clientSecret
	}
	if *ClientSecret == "" {
		return ErrClientSecretRequired
	}

	// MemcacheAddr
	if MemcacheAddr == nil || *MemcacheAddr == "" {
		memcacheAddr := os.Getenv(envPrefix + "MEMCACHE_ADDR")
		MemcacheAddr = &memcacheAddr
	}
	if *MemcacheAddr == "" {
		memcacheAddr := defaultMemcacheAddr
		MemcacheAddr = &memcacheAddr
	}

	return nil
}
