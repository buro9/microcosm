package ui

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Profile struct {
	ID             int64        `json:"id"`
	SiteID         int64        `json:"siteId,omitempty"`
	UserID         int64        `json:"userId"`
	Email          string       `json:"email,omitempty"`
	ProfileName    string       `json:"profileName"`
	Member         *bool        `json:"member,omitempty"`
	Gender         string       `json:"gender,omitempty"`
	Visible        *bool        `json:"visible"`
	StyleID        int64        `json:"styleId"`
	ItemCount      int32        `json:"itemCount"`
	CommentCount   int32        `json:"commentCount"`
	ProfileComment interface{}  `json:"profileComment"`
	Created        time.Time    `json:"created"`
	LastActive     time.Time    `json:"lastActive"`
	AvatarURL      string       `json:"avatar"`
	Meta           ExtendedMeta `json:"meta"`
}

type ProfileSummary struct {
	ID          int64        `json:"id"`
	SiteID      int64        `json:"siteId,omitempty"`
	UserID      int64        `json:"userId"`
	ProfileName string       `json:"profileName"`
	Visible     *bool        `json:"visible"`
	AvatarURL   string       `json:"avatar"`
	Meta        ExtendedMeta `json:"meta"`
}

type ProfileResponse struct {
	BoilerPlate
	Data Profile `json:"data"`
}

// userFromAPIContext is used to return the current user of the site based on
// the /whoami call determined by the apiroot that is within the context and the
// access token also in the context
func userFromAPIContext(ctx context.Context) (*Profile, error) {
	at := accessTokenFromContext(ctx)
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

	var apiResp ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return &apiResp.Data, nil
}
