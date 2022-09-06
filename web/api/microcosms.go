package api

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/buro9/microcosm/models"
)

// GetMicrocosm returns a microcosm if this user (defined by context) has
// permission to view it
func GetMicrocosm(ctx context.Context, microcosmID int64) (*models.Microcosm, error) {
	resp, err := apiGet(Params{Ctx: ctx, Type: "microcosms", TypeID: strconv.FormatInt(microcosmID, 10)})
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
