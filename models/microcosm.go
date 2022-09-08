package models

import "time"

// Microcosm describes a microcosm
type Microcosm struct {
	ID         int64 `json:"id"`
	SiteID     int64 `json:"siteId"`
	ParentID   int64 `json:"parentId,omitempty"`
	Breadcrumb []struct {
		Rel      string `json:"rel,omitempty"` // REST
		Href     string `json:"href"`
		Title    string `json:"title,omitempty"`
		Text     string `json:"text,omitempty"` // HTML
		LogoURL  string `json:"logoUrl,omitempty"`
		ID       int64  `json:"id"`
		Level    int64  `json:"level,omitempty"`
		ParentID int64  `json:"parentId,omitempty"`
	} `json:"breadcrumb,omitempty"`
	Visibility  string   `json:"visibility"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	LogoURL     string   `json:"logoUrl"`
	ItemTypes   []string `json:"itemTypes"`
	Moderators  []int64  `json:"moderators"`

	Items Array `json:"items"`
	Meta  struct {
		Created    time.Time       `json:"created"`
		CreatedBy  ProfileSummary  `json:"createdBy"`
		Edited     *time.Time      `json:"edited,omitempty"`
		EditedBy   *ProfileSummary `json:"editedBy,omitempty"`
		EditReason *string         `json:"editReason,omitempty"`
		Flags      struct {
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
		Stats       []Stat     `json:"stats,omitempty"`
		Links       []Link     `json:"links,omitempty"`
		Permissions Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// MicrocosmSummary describes the summary of a microcosm
type MicrocosmSummary struct {
	ID         int64 `json:"id"`
	ParentID   int64 `json:"parentId"`
	SiteID     int64 `json:"siteId"`
	Breadcrumb []struct {
		Rel      string `json:"rel,omitempty"` // REST
		Href     string `json:"href"`
		Title    string `json:"title,omitempty"`
		Text     string `json:"text,omitempty"` // HTML
		LogoURL  string `json:"logoUrl,omitempty"`
		ID       int64  `json:"id"`
		Level    int64  `json:"level,omitempty"`
		ParentID int64  `json:"parentId,omitempty"`
	} `json:"breadcrumb"`
	Visibility  string   `json:"visibility"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	LogoURL     string   `json:"logoUrl"`
	ItemTypes   []string `json:"itemTypes"`

	Children []struct {
		Rel      string `json:"rel,omitempty"` // REST
		Href     string `json:"href"`
		Title    string `json:"title,omitempty"`
		Text     string `json:"text,omitempty"` // HTML
		LogoURL  string `json:"logoUrl,omitempty"`
		ID       int64  `json:"id"`
		Level    int64  `json:"level,omitempty"`
		ParentID int64  `json:"parentId,omitempty"`
	} `json:"children"`
	Moderators   []int64 `json:"moderators"`
	ItemCount    int64   `json:"totalItems"`
	CommentCount int64   `json:"totalComments"`

	MostRecentUpdate SummaryItem `json:"mostRecentUpdate"`

	Meta struct {
		Created   time.Time      `json:"created"`
		CreatedBy ProfileSummary `json:"createdBy"`
		Flags     struct {
			Sticky    bool `json:"sticky"`
			Open      bool `json:"open"`
			Deleted   bool `json:"deleted"`
			Moderated bool `json:"moderated"`
			Visible   bool `json:"visible"`
			Unread    bool `json:"unread"`
			Watched   bool `json:"watched"`
			Ignored   bool `json:"ignored"`
			SendEmail bool `json:"sendEmail"`
			SendSMS   bool `json:"sendSMS"`
		} `json:"flags,omitempty"`
		Stats       []Stat      `json:"stats,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// MicrocosmResponse describes the API response that contains a microcosm
type MicrocosmResponse struct {
	BoilerPlate
	Data Microcosm `json:"data"`
}
