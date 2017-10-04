package models

import "encoding/json"

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
