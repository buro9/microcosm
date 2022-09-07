package controllers

import (
	"net/http"

	"github.com/buro9/microcosm/models"
	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
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

	conversation, status, err := api.GetConversation(r.Context(), conversationID, jumpTo, r.URL.Query())
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	// Horrible cludge of code to fetch attachments for any comments that has attachments
	var comments []models.Comment
	for _, comment := range *conversation.Items.AsComments() {
		if comment.Attachments > 0 {
			commentAttachments, status, err := api.GetCommentAttachments(r.Context(), comment.ID)
			if err != nil {
				errors.Render(w, r, status, err)
				return
			}
			comment.Files = *commentAttachments.Attachments.AsAttachments()
		}
		comments = append(comments,comment)
	}

	data := templates.Data{
		Request:    r,
		Site:       bag.GetSite(r.Context()),
		User:       bag.GetProfile(r.Context()),
		Section:    `home`,
		Pagination: models.ParsePagination(conversation.Items),

		Comments: &comments,
		Conversation: conversation,
	}

	err = templates.RenderHTML(w, "conversation", data)
	if err != nil {
		errors.Render(w, r, status, err)
	}
}
