package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/memcache"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/opts"
)

const (
	apiVersion  string = "/api/v1"
	apiMemcache string = "localhost:11211"
)

var (
	cnameToAPIRoot     map[string]string
	cnameToAPIRootLock sync.RWMutex
)

var apiCache = memcache.New(apiMemcache)

// ApiRootFromRequest returns the URL of the API for the site associated with
// the request, i.e. https://subdomain.apidomain.tld/api/v1
func ApiRootFromRequest(req *http.Request) (string, error) {
	// TODO: Remove this bit that forces LFGSS onto the world
	if time.Now().Year() == 2016 {
		return "https://lfgss.microco.sm" + apiVersion, nil
	}

	if strings.HasSuffix(req.Host, *opts.ApiDomain) {
		return "https://" + req.Host + apiVersion, nil
	}

	// Check cache
	cnameToAPIRootLock.RLock()
	apiURL, ok := cnameToAPIRoot[req.Host]
	cnameToAPIRootLock.RUnlock()
	if ok {
		return apiURL, nil
	}

	// Unknown host (custom name) so we need to ask the authority API service if
	// it knows of this site, and if so what the API domain is.
	resp, err := http.Get(
		fmt.Sprintf(
			"https://%s/api/v1/hosts/%s",
			*opts.ApiDomain,
			req.Host,
		),
	)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "",
			fmt.Errorf(
				"%s lookup failed: %s",
				req.Host,
				http.StatusText(resp.StatusCode),
			)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// the body should just be a text string of the host name which we can test
	// because it will end in the *opts.ApiDomain
	host := string(b)
	if strings.HasSuffix(host, *opts.ApiDomain) {
		u, err := url.Parse("https://" + string(b) + apiVersion)
		if err != nil {
			// extremely unlikely but this is insurance to check it was valid
			return "", err
		}

		// Add to cache
		cnameToAPIRootLock.Lock()
		cnameToAPIRoot[req.Host] = u.String()
		cnameToAPIRootLock.Unlock()
		return u.String(), nil
	}

	return "", fmt.Errorf("%s is not a valid host", host)
}

func buildAPIURL(ctx context.Context, endpoint string, q *url.Values) *url.URL {
	// ensure that we start with a trailing slash
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	// not possible to generate an error at this point, as this is not called
	// before newContext which already barfs on any failure to set the apiRoot
	u, _ := url.Parse(bag.GetAPIRoot(ctx) + endpoint)

	// Add any querystring args that were set
	if q != nil {
		u.RawQuery = q.Encode()
	}

	return u
}

// apiGet will perform an API call and if an error has occurred the body will
// have been read. If no error occurred the callee must read the body and ensure
// it is closed.
func apiGet(
	ctx context.Context,
	endpoint string,
	q *url.Values,
) (*http.Response, error) {

	u := buildAPIURL(ctx, endpoint, q)

	// Standard client using the cache transport
	c := &http.Client{
		Transport: httpcache.NewTransport(apiCache),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("User-Agent", "microcosm-ui")

	// Add auth if we have it, though we never use it for the "site" endpoint as
	// that is a perma-cache item
	if endpoint != "site" {
		if at := bag.GetAccessToken(ctx); at != "" {
			req.Header.Add("Authorization", "Bearer "+at)
		}
	}

	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}
	log.Printf("%s %s", u.String(), time.Since(start))

	return resp, errFromResp(resp)
}

// apiPost will perform an API call and if an error has occurred the body will
// have been read. If no error occurred the callee must read the body and ensure
// it is closed.
func apiPost(
	ctx context.Context,
	endpoint string,
	q *url.Values,
	data string,
) (*http.Response, error) {

	u := buildAPIURL(ctx, endpoint, q)

	c := &http.Client{}

	var br *bytes.Reader
	if data != "" {
		br = bytes.NewReader([]byte(data))
	}
	req, err := http.NewRequest("POST", u.String(), br)

	req.Header.Add("User-Agent", "microcosm-ui")
	if at := bag.GetAccessToken(ctx); at != "" {
		req.Header.Add("Authorization", "Bearer "+at)
	}

	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, errFromResp(resp)
}

func errFromResp(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		// Read the error and return that as the error
		defer resp.Body.Close()

		var boiler models.BoilerPlate
		err := json.NewDecoder(resp.Body).Decode(&boiler)
		if err != nil {
			return err
		}

		return fmt.Errorf(strings.Join(boiler.Errors, ","))
	}

	return nil
}
