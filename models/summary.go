package models

import (
	"encoding/json"
	"time"
)

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

	// Search results may come with a highlighted portion that matches the
	// searched term
	Highlight string `json:"highlight"`
}

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

func (m *SummaryItem) AsCommentSummary(raw json.RawMessage) *CommentSummary {
	if raw == nil {
		return nil
	}

	var summary CommentSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}

	return &summary
}

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
