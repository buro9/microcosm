package api

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/buro9/microcosm/models"
)

// GetConversation returns a conversation if this user (defined by context) has
// permission to view it
func GetConversation(ctx context.Context, conversationID int64, jumpTo string) (*models.Conversation, error) {
	resp, err := apiGet(Params{Ctx: ctx, Type: "conversations", TypeID: strconv.FormatInt(conversationID, 10), Part: jumpTo})
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
