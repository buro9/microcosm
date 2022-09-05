package models

import "encoding/json"

// Array describes arrays returned by the API
type Array struct {
	Total     int64           `json:"total"`
	Limit     int64           `json:"limit"`
	Offset    int64           `json:"offset"`
	MaxOffset int64           `json:"maxOffset"`
	Pages     int64           `json:"totalPages"`
	Page      int64           `json:"page"`
	Links     []Link          `json:"links,omitempty"`
	Type      string          `json:"type"`
	Items     json.RawMessage `json:"items"`
}

// AsComments will return the array items as comments
func (m *Array) AsComments() *[]Comment {
	if m.Items == nil {
		return nil
	}

	var comments []Comment
	if err := json.Unmarshal(m.Items, &comments); err != nil {
		return nil
	}

	return &comments
}

// AsHuddleSummaries will return the array items as summaries (of huddle type)
func (m *Array) AsHuddleSummaries() *[]HuddleSummary {
	if m.Items == nil {
		return nil
	}

	var summaries []HuddleSummary
	if err := json.Unmarshal(m.Items, &summaries); err != nil {
		return nil
	}

	return &summaries
}

// AsProfileSummaries will return the array items as summaries (of profile type)
func (m *Array) AsProfileSummaries() *[]ProfileSummary {
	if m.Items == nil {
		return nil
	}

	var summaries []ProfileSummary
	if err := json.Unmarshal(m.Items, &summaries); err != nil {
		return nil
	}

	return &summaries
}

// AsSummaryItems will return the array items as summaries (of mixed types)
func (m *Array) AsSummaryItems() *[]SummaryItem {
	if m.Items == nil {
		return nil
	}

	var summaries []SummaryItem
	if err := json.Unmarshal(m.Items, &summaries); err != nil {
		return nil
	}

	return &summaries
}
