# microcosm

Microcosm is forum software written in Go.

Logically it is implemented as:

* Web UI within a Go server that is itself a client of the API
* API within a Go server that interacts with a PostgreSQL database
* PostgreSQL database


This repository will eventually produce three binaries:

1. `microcosm-web` which is the web client and can be deployed and pointed at a JSON API
2. `microcosm-api` which is the JSON API server and talks to the database
3. `microcosm` which is a combined web application and API in one binary for single server or homogenous server installs

NOTE: Right now only `microcosm-web` is being produced, when the web UI is complete the existing [https://github.com/microcosm-cc/microcosm](https://github.com/microcosm-cc/microcosm) repo that contains the API will be merged or copied in once `microcosm-web` is complete as a standalone client.

All of the binaries can be load balanced and will scale horizontally. Only the database will hold state.

## Features

* Go web and API servers
* Within the Go web server all API responses are cached if possible (including all requests for guest users)
* Serves both HTTP and HTTPS (including HTTP/2) and forces SSL when required
* It's fast
* I can read this and make sense of it, meaning I can maintain it in future or others could

## Dependencies

* Go (latest stable preferred)
* Memcached running on localhost:11211 (the web client uses a httpcache transport for the API)
* PostgreSQL (latest stable preferred) with postgresql-contrib and the ltree module installed

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
  -keyFile string
      path to the TLS private key file (default "/etc/ssl/private/microco.sm.key")
  -listen string
      addr:port on which to serve HTTP (default ":80")
  -tlsListen string
      addr:port on which to serve HTTPS (default ":443")
```

The cert and key I'm pointing to represent a wildcard cert, you'll need one too if you want to serve more than a single website. Otherwise a Lets Encrypt cert is good enough.

Running the daemon is as simple as `make && sudo microcosm-web`. Why `sudo`? Because port 80 and 443... if you want to bind to those ports you need to sudo. Feel free to change the ports using the flags, and then place the front-end behind an nginx or similar to serve via 80 and 443. However... I intend to expose the UI directly to the world so that we take advantage of the strong default crypto in the Go web server, and for HTTP/2 out of the box.
