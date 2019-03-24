package models

import "time"

// DefaultMeta describes the universal metadata that all items have
type DefaultMeta struct {
	Created    time.Time       `json:"created"`
	CreatedBy  ProfileSummary  `json:"createdBy"`
	Edited     *time.Time      `json:"edited,omitempty"`
	EditedBy   *ProfileSummary `json:"editedBy,omitempty"`
	EditReason *string         `json:"editReason,omitempty"`
	ExtendedMeta
}

// SummaryMeta describes the metadata that summary items have
type SummaryMeta struct {
	Created   time.Time      `json:"created"`
	CreatedBy ProfileSummary `json:"createdBy"`
	ExtendedMeta
}

// ExtendedMeta is shared by DefaultMeta and SummaryMeta
type ExtendedMeta struct {
	Flags Flags `json:"flags,omitempty"`
	CoreMeta
}

// CoreMeta describes per item stats, links and permissions
type CoreMeta struct {
	Stats       []Stat      `json:"stats,omitempty"`
	Links       []Link      `json:"links,omitempty"`
	Permissions *Permission `json:"permissions,omitempty"`
}

// Flags describes per item flags
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
