# microcosm-ui
Front end for Microcosm, a Go web server that serves the static files, templates and performs API calls.

## Features

* Go
* API responses are cached if possible (notably for guest users)
* Serves both HTTP and HTTPS (including HTTP/2) and forces SSL when it's safe to
* It's fast
* I can read this and make sense of it

## Dependencies

* Go1.7+
* Memcached running on localhost:11211 (the web client uses a httpcache transport for the API)
* PostgreSQL 9.4+ with postgresql-contrib and the ltree module installed

NB: I'm currently not vendoring and am using the latest of a few different things. I will add vendoring when things have settled and the list of packages we're using has stabilised. You'll need to `go get` a few things to get things running.

## Usage

Use Lets Encrypt or buy a cert... you need a cert and key to run the UI as we'll force SSL on subdomain sites and any CNAME'd site that has `*Site.ForceSSL = true`.

flags reveal usage:

```
$ microcosm-web --help
Usage of microcosm-web:
  -apiDomain string
      the .tld that serves the API (default "microco.sm")
  -certFile string
      path to the TLS certificate file (default "/etc/ssl/certs/microco.sm.crt")
  -clientSecret string
      the API client secret (default os.Getenv("MICROCOSM_API_CLIENT_SECRET"))
  -files string
      directory that contains the templates and static files (default "/srv/microcosm-web")
  -keyFile string
      path to the TLS private key file (default "/etc/ssl/private/microco.sm.key")
  -listen string
      addr:port on which to serve HTTP (default ":80")
  -tlsListen string
      addr:port on which to serve HTTPS (default ":443")
```

I've symlinked `/srv/microcosm-ui` to point to `/home/buro9/Dev/src/github.com/microcosm-cc/microcosm-ui/files/` which is my local dev location for this repo. I expect to deploy to a prod environment where the static files actually are within `/srv/microcosm-ui`. You don't need to do this... you can just set the `-filespath``flag.

The cert and key I'm pointing to represent a wildcard cert, you'll need one too if you want to serve more than a single website. Otherwise a Lets Encrypt cert is good enough.

Running the daemon is as simple as `sudo microcosm-ui`. Why `sudo`? Port 80 and 443... if you want to bind to those ports you need to sudo. Feel free to change the ports using the flags, and then place the front-end behind an nginx or similar to serve via 80 and 443. However... I intend to expose the UI directly to the world so that we take advantage of the strong default crypto in the Go web server, and for HTTP/2 out of the box.
