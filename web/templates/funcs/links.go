package funcs

import "github.com/buro9/microcosm/models"

// LinkFromLinks will return a link given the rel="" value for it
func LinkFromLinks(links []models.Link, rel string) *models.Link {
	for _, link := range links {
		if link.Rel == rel {
			return &link
		}
	}

	return nil
}

// ReverseLinks will reverse a slice of links
func ReverseLinks(links []models.Link) []models.Link {
	var reversed []models.Link
	for i := len(links) - 1; i >= 0; i-- {
		reversed = append(reversed, links[i])
	}
	return reversed
}
