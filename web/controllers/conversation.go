package controllers

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

// ConversationGet will fetch a conversation
func ConversationGet(w http.ResponseWriter, r *http.Request) {
	conversationID := asInt64(r, "conversationID")
	var jumpTo string
	switch asString(r, "jumpTo") {
	case "newest":
		jumpTo = "newcomment"
	default:
	}

	conversation, err := api.GetConversation(r.Context(), conversationID, jumpTo)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(conversation.Items),

		Conversation: conversation,
	}

	err = templates.RenderHTML(w, "conversation", data)
	if err != nil {
		fmt.Printf("could not render %s\n", r.URL)
		w.Write([]byte(err.Error()))
	}
}
