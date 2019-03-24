package models

// UpdatesResponse descrives the API response for an array of updates
type UpdatesResponse struct {
	BoilerPlate
	Data UpdatesResults `json:"data"`
}

// UpdatesResults describes an array of updates
type UpdatesResults struct {
	Items Array `json:"updates"`
}
