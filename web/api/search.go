package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/buro9/microcosm/models"
)

func DoSearch(ctx context.Context, q url.Values) (*models.SearchResults, error) {
	resp, err := apiGet(ctx, "search", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.SearchResultsResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
