package models

import (
	"encoding/json"
	"time"
)

// SummaryItem describes the summary of an item
type SummaryItem struct {
	// The item type being summarised
	ItemType string `json:"itemType"`
	// The summary item
	Item json.RawMessage `json:"item"`

	// If this is an update, what was the trigger
	UpdateType string `json:"updateType"`
	// If this is an update, the identifier of this update
	ID int64 `json:"id"`

	// If this item has a parent (i.e. is a comment that belongs to an event)
	ParentItemType string `json:"parentItemType"`
	// A summary of the parent item
	ParentItem json.RawMessage `json:"parentItem"`

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
		} `json:"flags,omitempty"`
	} `json:"meta"`

	// Related to search results where the Item is not yet parsed and the page
	// needs to be able to consistently render something
	Unread       bool      `json:"unread"`
	Rank         float64   `json:"rank"`
	LastModified time.Time `json:"lastModified`
	Highlight    string    `json:"highlight"`
}

// AsItemSummary will return the raw input as a SummaryItem
func (m *SummaryItem) AsItemSummary(raw json.RawMessage) *ItemSummary {
	if raw == nil {
		return nil
	}

	var summary ItemSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsComment will return the raw input as a Comment
func (m *SummaryItem) AsCommentSummary(raw json.RawMessage) *Comment {
	if raw == nil {
		return nil
	}

	var summary Comment
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsConversationSummary will return the raw input as a ConversationSummary
func (m *SummaryItem) AsConversationSummary(raw json.RawMessage) *ConversationSummary {
	if raw == nil {
		return nil
	}

	var summary ConversationSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsEventSummary will return the raw input as a EventSummary
func (m *SummaryItem) AsEventSummary(raw json.RawMessage) *EventSummary {
	if raw == nil {
		return nil
	}

	var summary EventSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsHuddleSummary will return the raw input as a HuddleSummary
func (m *SummaryItem) AsHuddleSummary(raw json.RawMessage) *HuddleSummary {
	if raw == nil {
		return nil
	}

	var summary HuddleSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsMicrocosmSummary will return the raw input as a MicrocosmSummary
func (m *SummaryItem) AsMicrocosmSummary(raw json.RawMessage) *MicrocosmSummary {
	if raw == nil {
		return nil
	}

	var summary MicrocosmSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

// AsProfileSummary will return the raw input as a ProfileSummary
func (m *SummaryItem) AsProfileSummary(raw json.RawMessage) *ProfileSummary {
	if raw == nil {
		return nil
	}

	var summary ProfileSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}
