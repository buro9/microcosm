package models

import "time"

// // HuddleSummary describes a huddle summary returned by the API
// type HuddleSummary struct {
// 	ID          int64 `json:"id"`
// 	MicrocosmID int64 `json:"microcosmId"`
// 	Breadcrumb  []struct {
// 		Rel      string `json:"rel,omitempty"` // REST
// 		Href     string `json:"href"`
// 		Title    string `json:"title,omitempty"`
// 		Text     string `json:"text,omitempty"` // HTML
// 		LogoURL  string `json:"logoUrl,omitempty"`
// 		ID       int64  `json:"id"`
// 		Level    int64  `json:"level,omitempty"`
// 		ParentID int64  `json:"parentId,omitempty"`
// 	} `json:"breadcrumb,omitempty"`
// 	Title string `json:"title"`

// 	Participants []ProfileSummary `json:"participants"`

// 	CommentCount int64 `json:"totalComments"`
// 	ViewCount    int64 `json:"totalViews"`
// 	LastComment  struct {
// 		ID        int64          `json:"id"`
// 		Created   time.Time      `json:"created"`
// 		CreatedBy ProfileSummary `json:"createdBy"`
// 	} `json:"lastComment"`
// 	Meta struct {
// 		Created   time.Time      `json:"created"`
// 		CreatedBy ProfileSummary `json:"createdBy"`
// 		Flags     struct {
// 			Sticky    bool `json:"sticky,omitempty"`
// 			Open      bool `json:"open,omitempty"`
// 			Deleted   bool `json:"deleted,omitempty"`
// 			Moderated bool `json:"moderated,omitempty"`
// 			Visible   bool `json:"visible,omitempty"`
// 			Unread    bool `json:"unread,omitempty"`
// 			Watched   bool `json:"watched,omitempty"`
// 			Ignored   bool `json:"ignored,omitempty"`
// 			SendEmail bool `json:"sendEmail,omitempty"`
// 			SendSMS   bool `json:"sendSMS,omitempty"`
// 		} `json:"flags,omitempty"`
// 		Links       []Link      `json:"links,omitempty"`
// 		Permissions *Permission `json:"permissions,omitempty"`
// 	} `json:"meta"`
// }

// HuddleSummary is a summary of a huddle
type HuddleSummary struct {
	ID     int64  `json:"id"`
	SiteID int64  `json:"siteId,omitempty"`
	Title  string `json:"title"`

	CommentCount int64 `json:"totalComments"`

	Participants []ProfileSummary `json:"participants"`

	LastCommentID        int64          `json:"lastCommentId,omitempty"`
	LastCommentCreatedBy ProfileSummary `json:"lastCommentCreatedBy,omitempty"`
	LastCommentCreated   time.Time      `json:"lastCommentCreated,omitempty"`

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
			Attending bool `json:"attending,omitempty"`
		} `json:"flags,omitempty"`
		Stats       []Stat      `json:"stats,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// Huddles describes an array of huddles
type Huddles struct {
	Items Array `json:"huddles"`
	Meta  struct {
		Stats       []Stat     `json:"stats,omitempty"`
		Links       []Link     `json:"links,omitempty"`
		Permissions Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// HuddlesResponse describes the API response that contains an array of
// huddles
type HuddlesResponse struct {
	BoilerPlate
	Data Huddles `json:"data"`
}
