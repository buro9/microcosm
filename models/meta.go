package models

import "time"

type DefaultMeta struct {
	Created    time.Time       `json:"created"`
	CreatedBy  ProfileSummary  `json:"createdBy"`
	Edited     *time.Time      `json:"edited,omitempty"`
	EditedBy   *ProfileSummary `json:"editedBy,omitempty"`
	EditReason *string         `json:"editReason,omitempty"`
	ExtendedMeta
}

type SummaryMeta struct {
	Created   time.Time      `json:"created"`
	CreatedBy ProfileSummary `json:"createdBy"`
	ExtendedMeta
}

type ExtendedMeta struct {
	Flags Flags `json:"flags,omitempty"`
	CoreMeta
}

type CoreMeta struct {
	Stats       []Stat      `json:"stats,omitempty"`
	Links       []Link      `json:"links,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
}

type Flags struct {
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
}
