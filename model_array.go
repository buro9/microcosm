package ui

type Array struct {
	Total     int64         `json:"total"`
	Limit     int64         `json:"limit"`
	Offset    int64         `json:"offset"`
	MaxOffset int64         `json:"maxOffset"`
	Pages     int64         `json:"totalPages"`
	Page      int64         `json:"page"`
	Links     []Link        `json:"links,omitempty"`
	Type      string        `json:"type"`
	Items     []SummaryItem `json:"items"`
}
