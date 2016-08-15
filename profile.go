package ui

import "time"

type Profile struct {
	ID             int64        `json:"id"`
	SiteID         int64        `json:"siteId,omitempty"`
	UserID         int64        `json:"userId"`
	Email          string       `json:"email,omitempty"`
	ProfileName    string       `json:"profileName"`
	Member         *bool        `json:"member,omitempty"`
	Gender         string       `json:"gender,omitempty"`
	Visible        *bool        `json:"visible"`
	StyleID        int64        `json:"styleId"`
	ItemCount      int32        `json:"itemCount"`
	CommentCount   int32        `json:"commentCount"`
	ProfileComment interface{}  `json:"profileComment"`
	Created        time.Time    `json:"created"`
	LastActive     time.Time    `json:"lastActive"`
	AvatarURL      string       `json:"avatar"`
	Meta           ExtendedMeta `json:"meta"`
}

type ProfileSummary struct {
	ID          int64        `json:"id"`
	SiteID      int64        `json:"siteId,omitempty"`
	UserID      int64        `json:"userId"`
	ProfileName string       `json:"profileName"`
	Visible     *bool        `json:"visible"`
	AvatarURL   string       `json:"avatar"`
	Meta        ExtendedMeta `json:"meta"`
}
