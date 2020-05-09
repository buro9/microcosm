package models

import "time"

// Conversation describes a conversation returned by the API
type Conversation struct {
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

	Items Array `json:"comments"`

	Meta struct {
		Created   time.Time      `json:"created"`
		CreatedBy ProfileSummary `json:"createdBy"`
		Flags     struct {
			Sticky    bool `json:"sticky,omitempty"`
			Open      bool `json:"open,omitempty"`
			Deleted   bool `json:"deleted,omitempty"`
			Moderated bool `json:"moderated,omitempty"`
			Visible   bool `json:"visible,omitempty"`
			Unread    bool `json:"unread,omitempty"`
			Watched   bool `json:"watched,omitempty"`
			Ignored   bool `json:"ignored,omitempty"`
			SendEmail bool `json:"sendEmail,omitempty"`
			SendSMS   bool `json:"sendSMS,omitempty"`
		} `json:"flags,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// ConversationResponse describes the API response that contains a conversation
type ConversationResponse struct {
	BoilerPlate
	Data Conversation `json:"data"`
}

// ConversationSummary describes a conversation summary returned by the API
type ConversationSummary struct {
	ID          int64 `json:"id"`
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
	Title string `json:"title"`

	CommentCount int64 `json:"totalComments"`
	ViewCount    int64 `json:"totalViews"`
	LastComment  struct {
		ID        int64          `json:"id"`
		Created   time.Time      `json:"created"`
		CreatedBy ProfileSummary `json:"createdBy"`
	} `json:"lastComment"`
	Meta struct {
		Created   time.Time      `json:"created"`
		CreatedBy ProfileSummary `json:"createdBy"`
		Flags     struct {
			Sticky    bool `json:"sticky,omitempty"`
			Open      bool `json:"open,omitempty"`
			Deleted   bool `json:"deleted,omitempty"`
			Moderated bool `json:"moderated,omitempty"`
			Visible   bool `json:"visible,omitempty"`
			Unread    bool `json:"unread,omitempty"`
			Watched   bool `json:"watched,omitempty"`
			Ignored   bool `json:"ignored,omitempty"`
			SendEmail bool `json:"sendEmail,omitempty"`
			SendSMS   bool `json:"sendSMS,omitempty"`
		} `json:"flags,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}
