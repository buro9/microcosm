package templates

import "sync"

var loadDefinitionsOnce sync.Once

var page = []string{
	"breadcrumb",
	"pagination",
}

var microcosms = []string{
	"content_microcosm",
	"block_microcosm",
	"block_conversation",
	"block_event",
}

var searchResults = []string{
	"content_searchresults",
	"block_microcosm",
	"block_conversation",
	"block_event",
	"block_profile",
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
