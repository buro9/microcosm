package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/buro9/microcosm/models"
)

// GetUpdates returns the personalised list of items that have been updated for
// the given user (defined by context)
func GetUpdates(ctx context.Context, q url.Values) (*models.UpdatesResults, int, error) {
	resp, err := apiGet(Params{Ctx: ctx, Type: "updates", Q: q})
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	var apiResp models.UpdatesResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, resp.StatusCode, err
	}

	return &apiResp.Data, resp.StatusCode, nil
}
