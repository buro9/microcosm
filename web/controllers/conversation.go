package controllers

import (
	"fmt"
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/templates"
)

// ConversationGet will fetch the home page
func ConversationGet(w http.ResponseWriter, req *http.Request) {
	conversationID := asInt64(req, "conversationID")
	conversation, err := api.GetConversation(req.Context(), conversationID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	data := templates.Data{
		Request:    req,
		Site:       bag.GetSite(req.Context()),
		User:       bag.GetProfile(req.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(conversation.Items),

		Conversation: conversation,
	}

	err = templates.RenderHTML(w, "home", data)
	if err != nil {
		fmt.Println("could not render home")
		w.Write([]byte(err.Error()))
	}
}
