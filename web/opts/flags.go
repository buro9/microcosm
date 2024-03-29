package opts

import (
	"flag"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/securecookie"
)

const (
	envPrefix = "MICROCOSM_WEB_"

	defaultListen       = ":80"
	defaultTLSListen    = ":443"
	defaultCertFile     = "/etc/ssl/certs/microco.sm.crt"
	defaultKeyFile      = "/etc/ssl/private/microco.sm.key"
	defaultAPIDomain    = "microco.sm"
	defaultClientSecret = ""
	defaultMemcacheAddr = "localhost:11211"
	// TODO: Rotate and move to environment vars
	defaultCookieHashKey  = "70ce1fb50f865ef4f984fcb6fcabf1e870ce1fb50f865ef4f984fcb6fcabf1e8"
	defaultCookieBlockKey = "ed6f16535958f69087ccdd1556b6335d"

	defaultCsrfKey = "32-byte-long-auth-key"

	defaultIsDevelopment = false
)

var parseFlags sync.Once

var (
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

	// CookieHashKey is a 64 char string that is used to produce the HMAC of the
	// secure cookie used in the web app
	CookieHashKey *string

	// CookieBlockKey is a 32 char string that is used to encrypted the secure
	// cookie used in the web app
	CookieBlockKey *string

	// SecureCookie is an instance of gorilla securecookie that will be used
	// by the web app and underlines CSRF tokens
	SecureCookie *securecookie.SecureCookie

	// key to validate the CSRF token
	CsrfKey *string

	// Enables dev mode
	// Whether static assets are served from bundle
	IsDevelopment *bool
)

// RegisterFlags adds the flags needed by the UI if they have not already been
// added.
//
// Note: Defaults are empty as we will try os.Getenv for the values if they are
// provided by the command line, before finally setting defaults internally.
func RegisterFlags() {
	parseFlags.Do(
		func() {

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
	(default "`+defaultMemcacheAddr+`")`,
			)

			CookieHashKey = flag.String(
				"cookieHashKey",
				"",
				`the cookie HMAC is produced from this hash key (32 chars)
	alternatively $`+envPrefix+`COOKIE_HASH_KEY
	(default "`+defaultCookieHashKey+`")`,
			)

			CookieBlockKey = flag.String(
				"cookieBlockKey",
				"",
				`the cookie values are encrypted by this block key (64 chars)
	alternatively $`+envPrefix+`COOKIE_BLOCK_KEY
	(default "`+defaultCookieBlockKey+`")`,
			)

			CsrfKey = flag.String(
				"csrfKey",
				defaultCsrfKey,
				`the csrf token is validated with this key
	(default "`+defaultCsrfKey+`")`,
			)

			IsDevelopment = flag.Bool(
				"dev",
				defaultIsDevelopment,
				`is development
	(default "`+strconv.FormatBool(defaultIsDevelopment)+`")`,
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

	// CookieHashKey
	if CookieHashKey == nil || *CookieHashKey == "" {
		cookieHashKey := os.Getenv(envPrefix + "COOKIE_HASH_KEY")
		CookieHashKey = &cookieHashKey
	}
	if *CookieHashKey == "" {
		cookieHashKey := defaultCookieHashKey
		CookieHashKey = &cookieHashKey
	}

	// CookieBlockKey
	if CookieBlockKey == nil || *CookieBlockKey == "" {
		cookieBlockKey := os.Getenv(envPrefix + "COOKIE_BLOCK_KEY")
		CookieBlockKey = &cookieBlockKey
	}
	if *CookieBlockKey == "" {
		cookieBlockKey := defaultCookieBlockKey
		CookieBlockKey = &cookieBlockKey
	}

	// Create the instance of Secure Cookie that will be used during the life of
	// this program.
	SecureCookie = securecookie.New([]byte(*CookieHashKey), []byte(*CookieBlockKey))
	SecureCookie.MaxAge(60 * 60 * 24 * 365)

	return nil
}
