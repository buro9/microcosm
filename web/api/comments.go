package api

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/buro9/microcosm/models"
)

// GetCommentAttachments returns attachments for a given comment
func GetCommentAttachments(ctx context.Context, commentID int64) (*models.Attachments, error) {

	resp, err := apiGet(Params{Ctx: ctx, Type: "comments", TypeID: strconv.FormatInt(commentID, 10), Part: "attachments"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.AttachmentResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
