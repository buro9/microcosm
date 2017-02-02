package models

import "time"

type Profile struct {
	// Updateable
	ID          int64  `json:"id"`
	SiteID      int64  `json:"siteId,omitempty"`
	UserID      int64  `json:"userId"`
	Email       string `json:"email,omitempty"`
	ProfileName string `json:"profileName"`
	Member      bool   `json:"member,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Visible     bool   `json:"visible"`
	StyleID     int64  `json:"styleId"`
	AvatarURL   string `json:"avatar"`

	// Read only
	ItemCount      int64          `json:"itemCount"`
	CommentCount   int64          `json:"commentCount"`
	ProfileComment CommentSummary `json:"profileComment,omitempty"`

	Meta struct {
		Created    time.Time `json:"created"`
		LastActive time.Time `json:"lastActive"`
		Flags      struct {
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

type ProfileSummary struct {
	ID          int64        `json:"id"`
	SiteID      int64        `json:"siteId,omitempty"`
	UserID      int64        `json:"userId"`
	ProfileName string       `json:"profileName"`
	Visible     *bool        `json:"visible"`
	AvatarURL   string       `json:"avatar"`
	Meta        ExtendedMeta `json:"meta"`
}

type ProfileResponse struct {
	BoilerPlate
	Data Profile `json:"data"`
}
