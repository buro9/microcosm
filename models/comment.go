package models

import "time"

// Comment describes a comment summary returned by the API
type Comment struct {
	ID          int64        `json:"id"`
	ItemType    string       `json:"itemType"`
	ItemID      int64        `json:"itemId"`
	Revisions   int64        `json:"revisions"`
	InReplyTo   int64        `json:"inReplyTo"`
	Attachments int64        `json:"attachments"`
	FirstLine   string       `json:"firstLine"`
	Markdown    string       `json:"markdown"`
	HTML        string       `json:"html"`
	Files       []Attachment `json:"files,omitempty"`

	Meta struct {
		Created    time.Time       `json:"created"`
		CreatedBy  ProfileSummary  `json:"createdBy"`
		Edited     *time.Time      `json:"edited,omitempty"`
		EditedBy   *ProfileSummary `json:"editedBy,omitempty"`
		EditReason *string         `json:"editReason,omitempty"`
		Flags      struct {
			Deleted   bool `json:"deleted"`
			Moderated bool `json:"moderated"`
			Visible   bool `json:"visible"`
			Unread    bool `json:"unread"`
		} `json:"flags,omitempty"`
		Stats       []Stat      `json:"stats,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// CommentResponse describes the API response that contains a comment
type CommentResponse struct {
	BoilerPlate
	Data Comment `json:"data"`
}
