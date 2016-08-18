package ui

import "time"

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
