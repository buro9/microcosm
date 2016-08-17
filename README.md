# microcosm-ui
Front end for Microcosm, a Go web server that serves the static files, templates and performs API calls.

## Usage

Use Lets Encrypt or buy a cert... you need a cert and key to run the UI as we'll force SSL on subdomain sites and any CNAME'd site that has `*Site.ForceSSL = true`.

flags reveal usage:

```
$ microcosm-ui --help
Usage of microcosm-ui:
  -apiDomain string
    	the .tld that serves the API (default "microco.sm")
  -certFile string
    	path to the TLS certificate file (default "/etc/ssl/certs/microco.sm.crt")
  -filespath string
    	directory that contains the templates and static files (default "/srv/microcosm-ui")
  -keyFile string
    	path to the TLS private key file (default "/etc/ssl/private/microco.sm.key")
  -port int
    	port on which to serve HTTP (default 80)
  -tlsPort int
    	port on which to serve HTTPS (default 443)
```

I've symlinked `/srv/microcosm-ui` to point to `/home/buro9/Dev/src/github.com/microcosm-cc/microcosm-ui/files/` which is my local dev location for this repo. I expect to deploy to a prod environment where the static files actually are within `/srv/microcosm-ui`. You don't need to do this... you can just set the `-filespath``flag.

The cert and key I'm pointing to represent a wildcard cert, you'll need one too if you want to serve more than a single website. Otherwise a Lets Encrypt cert is good enough.

Running the daemon is as simple as `sudo microcosm-ui`. Why `sudo`? Port 80 and 443... if you want to bind to those ports you need to sudo. Feel free to change the ports using the flags, and then place the front-end behind an nginx or similar to serve via 80 and 443. However... I intend to expose the UI directly to the world so that we take advantage of the strong default crypto in the Go web server, and for HTTP/2 out of the box.
