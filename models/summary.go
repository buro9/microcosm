package models

import "encoding/json"

type SummaryItem struct {
	ItemType string          `json:"itemType"`
	Item     json.RawMessage `json:"item"`
}

func (item *SummaryItem) AsItemSummary() *ItemSummary {
	var m ItemSummary
	if err := json.Unmarshal(item.Item, &m); err != nil {
		return nil
	}

	return &m
}

func (item *SummaryItem) AsMicrocosmSummary() *MicrocosmSummary {
	var m MicrocosmSummary
	if err := json.Unmarshal(item.Item, &m); err != nil {
		return nil
	}

	return &m
}

func (item *SummaryItem) AsConversationSummary() *ConversationSummary {
	var m ConversationSummary
	if err := json.Unmarshal(item.Item, &m); err != nil {
		return nil
	}

	return &m
}

func (item *SummaryItem) AsEventSummary() *EventSummary {
	var m EventSummary
	if err := json.Unmarshal(item.Item, &m); err != nil {
		return nil
	}

	return &m
}
