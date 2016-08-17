package ui

import "time"

type ItemParent struct {
	MicrocosmID int64           `json:"microcosmId"`
	Breadcrumb  []MicrocosmLink `json:"breadcrumb,omitempty"`
}

type Item struct {
	ItemParent

	ItemType string `json:"itemType"`

	ID    int64  `json:"id"`
	Title string `json:"title"`

	CommentCount int64 `json:"totalComments"`
	ViewCount    int64 `json:"totalViews"`

	LastCommentID        int64          `json:"lastCommentId,omitempty"`
	LastCommentCreatedBy ProfileSummary `json:"lastCommentCreatedBy,omitempty"`
	LastCommentCreated   string         `json:"lastCommentCreated,omitempty"`

	Meta DefaultMeta `json:"meta"`
}

type ItemDetail struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	ItemParent
}

type ItemDetailCommentsAndMeta struct {
	Comments Array       `json:"comments"`
	Meta     DefaultMeta `json:"meta"`
}

type ItemSummary struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	ItemParent
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
