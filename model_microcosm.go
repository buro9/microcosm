package ui

import (
	"context"
	"encoding/json"
	"fmt"
)

type MicrocosmCore struct {
	ID          int64           `json:"id"`
	ParentID    int64           `json:"parentId,omitempty"`
	Breadcrumb  []MicrocosmLink `json:"breadcrumb,omitempty"`
	SiteID      int64           `json:"siteId,omitempty"`
	Visibility  string          `json:"visibility"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	LogoURL     string          `json:"logoUrl"`
	ItemTypes   []string        `json:"itemTypes"`
}

type Microcosm struct {
	MicrocosmCore

	RemoveLogo *bool `json:"removeLogo,omitempty"`

	Moderators []int64 `json:"moderators"`

	Items Array       `json:"items"`
	Meta  DefaultMeta `json:"meta"`
}

type MicrocosmSummary struct {
	MicrocosmCore

	Children     []MicrocosmLink `json:"children,omitempty"`
	Moderators   []int64         `json:"moderators"`
	ItemCount    int64           `json:"totalItems"`
	CommentCount int64           `json:"totalComments"`

	MostRecentUpdate *SummaryItem `json:"mostRecentUpdate,omitempty"`

	Meta SummaryMeta `json:"meta"`
}

type MicrocosmLink struct {
	Link

	LogoURL  string `json:"logoUrl,omitempty"`
	ID       int64  `json:"id"`
	Level    int64  `json:"level,omitempty"`
	ParentID int64  `json:"parentId,omitempty"`
}

type MicrocosmResponse struct {
	BoilerPlate
	Data Microcosm `json:"data"`
}

func microcosms(ctx context.Context) (*Microcosm, error) {
	resp, err := apiGet(ctx, "microcosms", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp MicrocosmResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &apiResp.Data, nil
}
