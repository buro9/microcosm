package api

import (
	"context"
	"encoding/json"
	"net/http"

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
