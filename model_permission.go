package ui

type Permission struct {
	CanCreate     bool `json:"create"`
	CanRead       bool `json:"read"`
	CanUpdate     bool `json:"update"`
	CanDelete     bool `json:"delete"`
	CanCloseOwn   bool `json:"closeOwn"`
	CanOpenOwn    bool `json:"openOwn"`
	CanReadOthers bool `json:"readOthers"`
	IsGuest       bool `json:"guest"`
	IsBanned      bool `json:"banned"`
	IsOwner       bool `json:"owner"`
	IsModerator   bool `json:"moderator"`
	IsSiteOwner   bool `json:"siteOwner"`
}
