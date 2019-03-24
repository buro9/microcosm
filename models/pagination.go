package models

// Pagination describes the pagination of an array
type Pagination struct {
	Page       int64  `json:"page"`
	TotalPages int64  `json:"totalPages"`
	Limit      int64  `json:"limit"`
	Offset     int64  `json:"offset"`
	MaxOffset  int64  `json:"maxOffset"`
	Links      []Link `json:"links"`
}

// ParsePagination returns the Pagination for an Array
func ParsePagination(array Array) *Pagination {
	p := Pagination{
		Page:       array.Page,
		TotalPages: array.Pages,
		Limit:      array.Limit,
		Offset:     array.Offset,
		Links:      array.Links,
	}

	return &p
}
