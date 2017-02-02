package models

type Pagination struct {
	Page       int64
	TotalPages int64
	Limit      int64
	Offset     int64
	Links      []Link
}

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
