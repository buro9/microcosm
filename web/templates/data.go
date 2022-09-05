package templates

import (
	"net/http"

	"github.com/buro9/microcosm/models"
)

// Data is the data that can be provided to a template.
//
// This is normalised into this one struct to ensure consistency across all
// templates, though very obviously not all templates require all fields and
// most of the time very few fields are filled in, typically an anonymous user
// will only have Site and whatever fields are relavent for a page shown, and a
// signed-in user will have Site and User along with whatever fields are
// relevant for the current page.
type Data struct {
	// Every request has these
	Request    *http.Request
	Site       *models.Site
	Section    string // Which part of the navbar to highlight
	Query      *models.SearchQuery
	Pagination *models.Pagination

	// If signed-in, this represents the signed-in user
	User *models.Profile

	// Depending on page and context, templates may expect the applicable one
	// to be filled in
	Array         *models.Array
	Conversation  *models.Conversation
	Huddles       *models.Huddles
	Microcosm     *models.Microcosm
	Profiles      *models.Profiles
	Profile       *models.Profile
	SearchResults *models.SearchResults
}
