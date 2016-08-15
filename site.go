package ui

import (
	"context"
	"encoding/json"
	"time"
)

type Site struct {
	ID                 int64          `json:"siteId"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	SubdomainKey       string         `json:"subdomainKey"`
	Domain             string         `json:"domain"`
	ForceSSL           *bool          `json:"forceSSL,omitempty"`
	OwnedBy            ProfileSummary `json:"ownedBy"`
	ThemeID            int64          `json:"themeId"`
	LogoURL            string         `json:"logoUrl"`
	FaviconURL         string         `json:"faviconUrl,omitempty"`
	BackgroundColor    string         `json:"backgroundColor"`
	BackgroundURL      string         `json:"backgroundUrl,omitempty"`
	BackgroundPosition string         `json:"backgroundPosition,omitempty"`
	LinkColor          string         `json:"linkColor"`
	GaWebPropertyID    string         `json:"gaWebPropertyId,omitempty"`
	Menu               []Link         `json:"menu"`

	Meta struct {
		Created    time.Time       `json:"created"`
		CreatedBy  ProfileSummary  `json:"createdBy"`
		Edited     *time.Time      `json:"edited,omitempty"`
		EditedBy   *ProfileSummary `json:"editedBy,omitempty"`
		EditReason *string         `json:"editReason,omitempty"`

		Flags struct {
			Deleted *bool `json:"deleted"`
		} `json:"flags,omitempty"`

		Stats       []Stat      `json:"stats,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		Permissions *Permission `json:"permissions,omitempty"`
	} `json:"meta"`
}

type SiteResponse struct {
	BoilerPlate
	Data Site `json:"data"`
}

// siteFromAPIContext is used to return a Site given the apiRoot that is within
// the context.
//
// TODO: If no access_token is present, cache the Site
func siteFromAPIContext(ctx context.Context) (*Site, error) {
	resp, err := apiGet(ctx, "site", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp SiteResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	site := &apiResp.Data
	if site.ThemeID == 0 {
		site.ThemeID = 1
	}

	return site, nil
}
