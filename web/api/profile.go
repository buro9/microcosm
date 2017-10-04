package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/bag"
)

// ProfileFromAPIContext is used to return the current profile of the site based
// on the /whoami call determined by the apiroot that is within the context and
// the access token also in the context
func ProfileFromAPIContext(ctx context.Context) (*models.Profile, error) {
	accessToken := bag.GetAccessToken(ctx)
	if accessToken == "" {
		// No access token means that we cannot possibly have a user
		return nil, nil
	}

	resp, err := apiGet(ctx, "whoami", nil)
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

// GetProfiles returns a list of profiles for a given search.
func GetProfiles(ctx context.Context, query url.Values) (*models.Profiles, error) {
	// Set the query options
	q := url.Values{}
	orderByComment := (query.Get("top") == strings.ToLower("true"))
	if orderByComment {
		q.Add("top", "true")
	}
	var nameStartsWith string
	if query.Get("q") != "" {
		v := strings.TrimLeft(query.Get("q"), "+@")
		if v != "" {
			nameStartsWith = v
			q.Add("q", v)
		}
	}
	isFollowing := (query.Get("following") == strings.ToLower("true"))
	if isFollowing {
		q.Add("following", "true")
	}
	isOnline := (query.Get("online") == strings.ToLower("true"))
	if isOnline {
		q.Add("online", "true")
	}
	offset := query.Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	resp, err := apiGet(ctx, "profiles", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.ProfilesResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	apiResp.Data.Query.Following = isFollowing
	apiResp.Data.Query.Online = isOnline
	apiResp.Data.Query.Q = nameStartsWith
	apiResp.Data.Query.Top = orderByComment

	return &apiResp.Data, nil
}
