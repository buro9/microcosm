package ui

type ConversationSummary struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	CommentCount         int64          `json:"totalComments"`
	ViewCount            int64          `json:"totalViews"`
	LastCommentID        int64          `json:"lastCommentId,omitempty"`
	LastCommentCreatedBy ProfileSummary `json:"lastCommentCreatedBy,omitempty"`
	LastCommentCreated   string         `json:"lastCommentCreated,omitempty"`

	MicrocosmID int64           `json:"microcosmId"`
	Breadcrumb  []MicrocosmLink `json:"breadcrumb,omitempty"`

	ItemSummaryMeta
}

type Conversation struct {
	ItemDetail
	ItemDetailCommentsAndMeta
}
