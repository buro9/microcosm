package api

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"

	"github.com/buro9/microcosm/models"
)

// GetConversation returns a conversation if this user (defined by context) has
// permission to view it
func GetConversation(ctx context.Context, conversationID int64, jumpTo string, query url.Values) (*models.Conversation, error) {

// Set the query options
	q := url.Values{}
	offset := query.Get("offset")
	if offset != "" {
		q.Add("offset", offset)
	}

	resp, err := apiGet(Params{Ctx: ctx, Type: "conversations", TypeID: strconv.FormatInt(conversationID, 10), Part: jumpTo, Q: q})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp models.ConversationResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &apiResp.Data, nil
}
