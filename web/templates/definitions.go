package templates

import "sync"

var loadDefinitionsOnce sync.Once

var page = []string{
	"breadcrumb",
	"pagination",
}

var microcosms = []string{
	"block_conversation",
	"block_event",
	"block_microcosm",
	"content_microcosm",
}

var profiles = []string{
	"block_profile",
	"content_profiles",
}

var profile = []string{
	"block_conversation",
	"block_event",
	"block_huddle",
	"block_list_comment",
	"block_microcosm",
	"content_profile",
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
				Template{
					Base:     "base",
					Page:     "home",
					Includes: Collate("sidebar_home", page, microcosms),
				},
				Template{
					Base:     "base",
					Page:     "microcosm",
					Includes: Collate("sidebar_microcosm", page, microcosms),
				},
				Template{
					Base:     "base",
					Page:     "profiles",
					Includes: Collate("sidebar_profiles", page, profiles),
				},
				Template{
					Base:     "base",
					Page:     "profile",
					Includes: Collate("sidebar_profile", page, profile),
				},
				Template{
					Base:     "base",
					Page:     "today",
					Includes: Collate("sidebar_today", page, searchResults),
				},
				Template{
					Base:     "base",
					Page:     "updates",
					Includes: Collate("sidebar_updates", page, searchResults),
				},
			}
		},
	)
}
