package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/bag"
)

// ProfileFromAPIContext is used to return the current profile of the site based
// on the /whoami call determined by the apiroot that is within the context and
// the access token also in the context
func ProfileFromAPIContext(ctx context.Context) (*models.Profile, error) {
	at := bag.GetAccessToken(ctx)
	if at == "" {
		// No access token means that we cannot possibly have a user
		return nil, nil
	}

	// This call will always perform a redirect and Go does not retain headers
	// on redirects, so we need to put the access token in the URL for this one
	// call and then totally bypass the apiGet caching layer (which would cache
	// the unauthenticated redirected call that wouldn't have the auth header).
	q := &url.Values{}
	q.Add("access_token", at)
	u := buildAPIURL(ctx, "whoami", q)

	c := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("User-Agent", "microcosm-ui")

	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("%s %s", u.String(), time.Since(start))

	if errFromResp(resp) != nil {
		return nil, err
	}

	if err != nil {
		switch resp.StatusCode {
		case http.StatusUnauthorized: // 401
			// Expired/invalid token, they are no longer signed-in
			return nil, nil
		case http.StatusForbidden: // 403
			// No access token provided
			return nil, nil
		case http.StatusNotFound: // 404
			// Valid access token for a now deleted or banned user
			return nil, nil
		default:
			// An unexpected error
			return nil, err
		}
	}
	defer resp.Body.Close()

	var apiResp models.ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp.Data, nil
}
