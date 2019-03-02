package api

import (
	"context"
	"encoding/json"
	"log"

	"github.com/buro9/microcosm/models"
)

// SiteFromAPIContext is used to return a Site given the apiRoot that is within
// the context.
func SiteFromAPIContext(ctx context.Context) (*models.Site, error) {
	resp, err := apiGet(ctx, "site", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.SiteResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
