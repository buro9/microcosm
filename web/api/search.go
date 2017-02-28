package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/buro9/microcosm/models"
)

func DoSearch(ctx context.Context, q *url.Values) (*models.SearchResults, error) {
	resp, err := apiGet(ctx, "search", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.SearchResultsResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &apiResp.Data, nil
}
