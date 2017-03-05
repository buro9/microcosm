package models

type UpdatesResponse struct {
	BoilerPlate
	Data UpdatesResults `json:"data"`
}

type UpdatesResults struct {
	Items Array `json:"updates"`
}
