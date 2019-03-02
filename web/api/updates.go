package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/buro9/microcosm/models"
)

func GetUpdates(ctx context.Context, q url.Values) (*models.UpdatesResults, error) {
	resp, err := apiGet(ctx, "updates", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.UpdatesResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
