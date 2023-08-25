package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gregjones/httpcache"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/opts"
)

const (
	apiVersion string = "/api/v1"
	userAgent  string = "microcosm Go web client"
)

var (
	cnameToAPIRoot     map[string]string
	cnameToAPIRootLock sync.RWMutex
)

// RootFromRequest returns the URL of the API for the site associated with
// the request, i.e. https://subdomain.apidomain.tld/api/v1
func RootFromRequest(req *http.Request) (string, int, error) {
	if strings.HasSuffix(req.Host, *opts.APIDomain) {
		return "https://" + req.Host + apiVersion, 0, nil
	}

	// Check cache
	cnameToAPIRootLock.RLock()
	apiURL, ok := cnameToAPIRoot[req.Host]
	cnameToAPIRootLock.RUnlock()
	if ok {
		return apiURL, 0, nil
	}

	// Unknown host (custom name) so we need to ask the authority API service if
	// it knows of this site, and if so what the API domain is.
	resp, err := http.Get(
		fmt.Sprintf(
			"https://%s/api/v1/hosts/%s",
			*opts.APIDomain,
			req.Host,
		),
	)
	if err != nil {
		log.Print(err)
		return "", resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode,
			fmt.Errorf(
				"%s lookup failed: %s",
				req.Host,
				http.StatusText(resp.StatusCode),
			)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return "", resp.StatusCode, err
	}

	// the body should just be a text string of the host name which we can test
	// because it will end in the *opts.ApiDomain
	host := string(b)
	if strings.HasSuffix(host, *opts.APIDomain) {
		u, err := url.Parse("https://" + string(b) + apiVersion)
		if err != nil {
			// extremely unlikely but this is insurance to check it was valid
			return "", resp.StatusCode, err
		}

		// Add to cache
		cnameToAPIRootLock.Lock()
		cnameToAPIRoot[req.Host] = u.String()
		cnameToAPIRootLock.Unlock()
		return u.String(), resp.StatusCode, nil
	}

	return "", resp.StatusCode, fmt.Errorf("%s is not a valid host", host)
}

// apiGet will perform an API call and if an error has occurred the body will
// have been read. If no error occurred the callee must read the body and ensure
// it is closed.
func apiGet(params Params) (*http.Response, error) {

	u := params.buildAPIURL()

	accessToken := bag.GetAccessToken(params.Ctx)

	var c *http.Client
	if (params.Type == "site" ||
		params.Type == "profiles" ||
		accessToken == "") && apiCache != nil {
		// Standard client using the cache transport for non-authenticated API
		// requests
		c = &http.Client{
			// CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 	return http.ErrUseLastResponse
			// },
			Transport: httpcache.NewTransport(apiCache),
		}
	} else {
		// Context cancellable transport for authenticated API requests
		c = &http.Client{
			// CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 	return http.ErrUseLastResponse
			// },
		}
	}

	req, _ := http.NewRequestWithContext(params.Ctx, "GET", u.String(), nil)
	req.Header.Add("User-Agent", userAgent)
	//req.Header.Add("X-Disable-Boiler", "true")

	// Add auth if we have it, though we never use it for the "site" endpoint as
	// that is a perma-cache item
	if accessToken != "" && (params.Type != "site" &&
		params.Type != "profiles" &&
		params.Type != "") {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}
	log.Printf(
		"[%s] GET %s %d %s",
		middleware.GetReqID(params.Ctx),
		u.String(),
		resp.StatusCode,
		time.Since(start),
	)

	return resp, errFromResp(resp)
}

// apiPost will perform an API call and if an error has occurred the body will
// have been read. If no error occurred the callee must read the body and ensure
// it is closed.
func apiPost(params Params, data interface{}) (*http.Response, error) {

	u := params.buildAPIURL()

	c := http.DefaultClient

	var br *bytes.Reader
	if data != nil {
		bs, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		br = bytes.NewReader(bs)
	}
	req, _ := http.NewRequest("POST", u.String(), br)

	req.Header.Add("User-Agent", userAgent)
	//req.Header.Add("X-Disable-Boiler", "true")
	req.Header.Add("Content-Type", "application/json")
	if at := bag.GetAccessToken(params.Ctx); at != "" {
		req.Header.Add("Authorization", "Bearer "+at)
	}

	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}
	log.Printf(
		"[%s] POST %s %d %s",
		middleware.GetReqID(params.Ctx),
		u.String(),
		resp.StatusCode,
		time.Since(start),
	)

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
