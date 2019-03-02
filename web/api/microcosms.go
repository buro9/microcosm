package api

import (
	"context"
	"encoding/json"
	"log"

	"github.com/buro9/microcosm/models"
)

func GetMicrocosms(ctx context.Context) (*models.Microcosm, error) {
	resp, err := apiGet(ctx, "microcosms", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.MicrocosmResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
