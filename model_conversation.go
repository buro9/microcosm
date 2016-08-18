package ui

type ConversationSummary struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	CommentCount         int64          `json:"totalComments"`
	ViewCount            int64          `json:"totalViews"`
	LastCommentID        int64          `json:"lastCommentId,omitempty"`
	LastCommentCreatedBy ProfileSummary `json:"lastCommentCreatedBy,omitempty"`
	LastCommentCreated   string         `json:"lastCommentCreated,omitempty"`

	MicrocosmID int64 `json:"microcosmId"`
	Breadcrumb  []struct {
		Rel      string `json:"rel,omitempty"` // REST
		Href     string `json:"href"`
		Title    string `json:"title,omitempty"`
		Text     string `json:"text,omitempty"` // HTML
		LogoURL  string `json:"logoUrl,omitempty"`
		ID       int64  `json:"id"`
		Level    int64  `json:"level,omitempty"`
		ParentID int64  `json:"parentId,omitempty"`
	} `json:"breadcrumb,omitempty"`

	ItemSummaryMeta
}

type Conversation struct {
	ItemDetail
	ItemDetailCommentsAndMeta
}
