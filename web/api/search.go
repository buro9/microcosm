package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/buro9/microcosm/models"
)

// DoSearch will perform a search against the search API for a given query
func DoSearch(ctx context.Context, q url.Values) (*models.SearchResults, error) {
	resp, err := apiGet(Params{Ctx: ctx, Type: "search", Q: q})
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
