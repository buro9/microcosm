package ui

import "time"

type Item struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	MicrocosmID int64  `json:"microcosmId"`
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

	ItemType string `json:"itemType"`

	CommentCount         int64          `json:"totalComments"`
	ViewCount            int64          `json:"totalViews"`
	LastCommentID        int64          `json:"lastCommentId,omitempty"`
	LastCommentCreatedBy ProfileSummary `json:"lastCommentCreatedBy,omitempty"`
	LastCommentCreated   string         `json:"lastCommentCreated,omitempty"`

	Meta DefaultMeta `json:"meta"`
}

type ItemDetail struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	MicrocosmID int64  `json:"microcosmId"`
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
}

type ItemDetailCommentsAndMeta struct {
	Comments Array       `json:"comments"`
	Meta     DefaultMeta `json:"meta"`
}

type ItemSummary struct {
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
}

type ItemSummaryMeta struct {
	CommentCount int64       `json:"totalComments"`
	ViewCount    int64       `json:"totalViews"`
	LastComment  LastComment `json:"lastComment,omitempty"`
	Meta         SummaryMeta `json:"meta"`
}

type LastComment struct {
	ID        int64          `json:"id"`
	Created   time.Time      `json:"created"`
	CreatedBy ProfileSummary `json:"createdBy"`
}
