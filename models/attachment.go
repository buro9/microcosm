package models

import "time"

// Attachment describes an attachment returned by the API
type Attachment struct {
	ProfileID int64     `json:"profileId"`
	FileHash  string    `json:"fileHash"`
	FileName  string    `json:"fileName"`
	FileExt   string    `json:"fileExt"`
	Created   time.Time `json:"created"`

	Meta struct {
		Stats       []Stat     `json:"stats,omitempty"`
		Links       []Link     `json:"links,omitempty"`
		Permissions Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

type Attachments struct {
	Attachments Array `json:attachments`
	Meta        struct {
		Links       []Link     `json:"links,omitempty"`
		Permissions Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// AttachmentResponse describes the API response that contains a conversation
type AttachmentResponse struct {
	BoilerPlate
	Data Attachments `json:"data"`
}
