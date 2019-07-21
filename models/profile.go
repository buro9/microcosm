package models

import "time"

// Profile describes a profile
type Profile struct {
	// Updateable
	ID          int64     `json:"id"`
	SiteID      int64     `json:"siteId,omitempty"`
	UserID      int64     `json:"userId"`
	Email       string    `json:"email,omitempty"`
	ProfileName string    `json:"profileName"`
	Member      bool      `json:"member,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	Visible     bool      `json:"visible"`
	StyleID     int64     `json:"styleId"`
	AvatarURL   string    `json:"avatar"`
	Created     time.Time `json:"created"`
	LastActive  time.Time `json:"lastActive"`

	// Read only
	ItemCount      int64          `json:"itemCount"`
	CommentCount   int64          `json:"commentCount"`
	ProfileComment CommentSummary `json:"profileComment,omitempty"`

	Meta struct {
		Flags struct {
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

// ProfileSummary summarises a profile
type ProfileSummary struct {
	ID          int64        `json:"id"`
	SiteID      int64        `json:"siteId,omitempty"`
	UserID      int64        `json:"userId"`
	ProfileName string       `json:"profileName"`
	Visible     *bool        `json:"visible"`
	AvatarURL   string       `json:"avatar"`
	Meta        ExtendedMeta `json:"meta"`
}

// ProfileResponse describes the API response that contains a Profile
type ProfileResponse struct {
	BoilerPlate
	Data Profile `json:"data"`
}

// Profiles describes an array of profiles
type Profiles struct {
	Query struct {
		Following bool   `json:"following,omitempty"`
		Online    bool   `json:"online,omitempty"`
		Q         string `json:"q"`
		Top       bool   `json:"top,omitempty"`
	} `json:"query"`
	Items Array `json:"profiles"`
	Meta  struct {
		Flags struct {
			Watched   bool `json:"watched,omitempty"`
			SendEmail bool `json:"sendEmail,omitempty"`
			SendSMS   bool `json:"sendSMS,omitempty"`
		} `json:"flags,omitempty"`
		Links       []Link     `json:"links,omitempty"`
		Permissions Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

// ProfilesResponse describes the API response that contains an array of
// profiles
type ProfilesResponse struct {
	BoilerPlate
	Data Profiles `json:"data"`
}
