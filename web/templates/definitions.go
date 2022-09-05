package templates

import "sync"

var loadDefinitionsOnce sync.Once

var conversation = []string{
	"block_list_comment",
	"content_conversation",
}

var huddles = []string{
	"block_huddle",
	"content_huddles",
}

var microcosms = []string{
	"block_conversation",
	"block_event",
	"block_microcosm",
	"content_microcosm",
}

var page = []string{
	"breadcrumb",
	"pagination",
}

var profile = []string{
	"block_conversation",
	"block_event",
	"block_huddle",
	"block_list_comment",
	"block_microcosm",
	"content_profile",
}

var profiles = []string{
	"block_profile",
	"content_profiles",
}

var searchResults = []string{
	"block_conversation",
	"block_event",
	"block_huddle",
	"block_list_comment",
	"block_microcosm",
	"block_profile",
	"content_searchresults",
}

func loadDefinitions() {
	loadDefinitionsOnce.Do(
		func() {
			Templates = []Template{
				{
					Base:     "base",
					Page:     "conversation",
					Includes: Collate("sidebar_conversation", page, conversation),
				},
				{
					Base:     "base",
					Page:     "home",
					Includes: Collate("sidebar_home", page, microcosms),
				},
				{
					Base:     "base",
					Page:     "huddles",
					Includes: Collate("sidebar_huddles", page, huddles),
				},
				{
					Base:     "base",
					Page:     "microcosm",
					Includes: Collate("sidebar_microcosm", page, microcosms),
				},
				{
					Base:     "base",
					Page:     "profile",
					Includes: Collate("sidebar_profile", page, profile),
				},
				{
					Base:     "base",
					Page:     "profiles",
					Includes: Collate("sidebar_profiles", page, profiles),
				},
				{
					Base:     "base",
					Page:     "today",
					Includes: Collate("sidebar_today", page, searchResults),
				},
				{
					Base:     "base",
					Page:     "updates",
					Includes: Collate("sidebar_updates", page, searchResults),
				},
			}
		},
	)
}
