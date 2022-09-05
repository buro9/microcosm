package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/buro9/microcosm/models"
)

// GetHuddles returns a list of huddles for a given search.
func GetHuddles(ctx context.Context, query url.Values) (*models.Huddles, error) {
	// Set the query options
	q := url.Values{}
	unread := (query.Get("unread") == strings.ToLower("true"))
	if unread {
		q.Add("unread", "true")
	}
	offset := query.Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	resp, err := apiGet(Params{Ctx: ctx, PathPrefix: "huddles", Q: q})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.HuddlesResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
