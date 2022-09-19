package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/buro9/microcosm/models"
)

// CommentForm is a comment as the form posts it.
//
// We define this here even though it's a valid subset of the Comment struct
// and we could've used that, this will give us more protection against
// attacks that targeted the unreferenced fields. In this, only the fields that
// a form could populate are provided, the we will perform the mapping to the
// models/comment.go within the controller, which is inaccessible to an attacker.
//
// In essence this is the write path, and non-Form structs are the read path
// and that abstraction protects us. It also allows us to define validation here.
type CommentForm struct {
	ID          int64  `json:"id"`
	ItemType    string `json:"itemType",validate:"required"`
	ItemID      int64  `json:"itemId",validate:"required"`
	InReplyTo   int64  `json:"inReplyTo"`
	Markdown    string `json:"markdown",validate:"required,max=50000"`
	Attachments int64  `json:"attachments"`
}

// GetCommentAttachments returns attachments for a given comment
func GetCommentAttachments(ctx context.Context, commentID int64) (*models.Attachments, int, error) {

	resp, err := apiGet(Params{Ctx: ctx, Type: "comments", TypeID: strconv.FormatInt(commentID, 10), Part: "attachments"})
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	var apiResp models.AttachmentResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		log.Print(err)
		return nil, resp.StatusCode, err
	}

	return &apiResp.Data, resp.StatusCode, nil
}

func (commentForm *CommentForm) Post(ctx context.Context) (*models.Comment, int, error) {
	resp, err := apiPost(
		Params{
			Ctx:  ctx,
			Type: "comments",
		},
		commentForm,
	)

	if err != nil {
		log.Printf("apiPost(`comment`): %s", err.Error())
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	// The POST has done an implicit round-trip... it has POST'd the commentForm
	// and received a 302 that points to the new comment, and then it has followed
	// the Location header and retrieved the comment. What exists in the body
	// is the comment... and within the meta struct for that we get the URL
	// to the comment.
	//
	// All we need to do is return the comment and have the callee handle sending
	// the frontend to the right place.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("apiPost(`comment`): %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	var commentResponse models.CommentResponse
	if err := json.Unmarshal(body, &commentResponse); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &commentResponse.Data, http.StatusOK, nil
}
