package api

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/buro9/microcosm/web/bag"
)

// Params is a deconstruction of a Microcosm API URL into constiuent parts that
// can be built up to form an API request.
//
// For examples of URLs, see the Microcosm API server handlers.
type Params struct {
	Ctx context.Context
	// Type corresponds to the top level API path.
	// i.e. comments, conversation, events, huddles, legal, microcosms, polls, profiles, roles, updates, users, watchers
	Type string
	// TypeID corresponds to a single item of the given Type
	TypeID string
	// Part corresponds to a child of the Type and TypeID
	// i.e. /comments/17/attachments
	// This would make Part = "attachments"
	Part string
	// PartID corresponds to a single item of the given Part for the given Type and TypeID
	// i.e. /comments/17/attachments/<filehash>
	// This would make PartID = "<filehash>" (assume a SHA file hash is the actual value)
	PartID string
	// Part2 corresponds to a child of the Part
	// i.e. /microcosms/31/roles/18/profiles
	// This would make Part2 = "profiles"
	Part2 string
	// Part2ID corresponds to a single item of the given Part2, for the Given Type + TypeID + Part + PartID
	// i.e. //microcosms/31/roles/18/profiles/68
	// This would make Part2ID = "68"
	Part2ID string
	// Q is the query string
	Q url.Values
}

func (p Params) buildAPIURL() *url.URL {
	// ensure that we start with a trailing slash
	if !strings.HasPrefix(p.Type, "/") {
		p.Type = "/" + p.Type
	}

	// It is not possible to generate an error at this point, as this is not called
	// before newContext which already barfs on any failure to set the apiRoot
	//
	// If tests are written for this func without calling newContext then an error
	// may be possible
	//
	// NB: params.Type is always added
	u, _ := url.Parse(bag.GetAPIRoot(p.Ctx) + p.Type)

	var path strings.Builder
	path.WriteString(u.Path)
	if p.TypeID != "" {
		fmt.Fprintf(&path, "/%s", p.TypeID)
	}
	if p.Part != "" {
		fmt.Fprintf(&path, "/%s", p.Part)
	}
	if p.PartID != "" {
		fmt.Fprintf(&path, "/%s", p.PartID)
	}
	if p.Part2 != "" {
		fmt.Fprintf(&path, "/%s", p.Part2)
	}
	if p.Part2ID != "" {
		fmt.Fprintf(&path, "/%s", p.Part2ID)
	}
	u.Path = path.String()

	// Add any querystring args that were set
	if p.Q != nil {
		u.RawQuery = p.Q.Encode()
	}

	return u
}
