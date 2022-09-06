package funcs

import (
	"fmt"
	net_url "net/url"
	"strings"
)

// api2ui will convert a relative Microcosm API URL into a URL for the
// equivalent item as viewed by the Web UI.
//
// For example:
// * API URL: /api/v1/conversations/13242
// * UI URL:  /conversations/13242/
//
// Simply:
// 1. Strip API prefix
// 2. API resources don't have trailing slashes but web pages in Microcosm do
//
// If an error is encountered, no string is returned... ensure your input is a
// URL.
func api2ui(s string) string {
	u, err := net_url.Parse(s)
	if err != nil {
		return ""
	}

	u.Path = strings.Replace(u.Path, `/api/v1`, ``, 1)

	if !strings.HasSuffix(u.Path, `/`) {
		u.Path += `/`
	}

	return u.String()
}

// url is not intelligent, it expects to be passed the right thing and will
// return nothing if the args are not correct for the given key
//
// This is effectively how all links in the front-end are constructed
func url(key string, args ...interface{}) string {
	switch len(args) {
	case 0:
		switch key {
		case "home":
			return "/"
		case "huddle-create":
			return "/huddles/create/"
		case "huddle-list":
			return "/huddles/"
		case "legal-list":
			return "/legal/"
		case "login":
			return "/login/"
		case "logout":
			return "/logout/"
		case "profile-list":
			return "/profiles/"
		case "search":
			return "/search/"
		case "today":
			return "/today/"
		case "update-list":
			return "/updates/"
		case "update-settings":
			return "/updates/settings/"
		case "watcher":
			return "/watchers/"
		}
	case 1:
		switch v := args[0].(type) {
		case string:
			switch key {
			case "legal":
				return fmt.Sprintf("/legal/%s/", v)
			}
		case int64:
			switch key {
			case "comment":
				return fmt.Sprintf("/comments/%d/", v)
			case "comment-incontext":
				return fmt.Sprintf("/comments/%d/incontext/", v)
			case "conversation-create":
				return fmt.Sprintf("/microcosms/%d/create/conversation/", v)
			case "event-create":
				return fmt.Sprintf("/microcosms/%d/create/event/", v)
			case "huddle":
				return fmt.Sprintf("/huddles/%d/", v)
			case "huddle-newest":
				return fmt.Sprintf("/huddles/%d/newest/", v)
			case "memberships-list":
				return fmt.Sprintf("/microcosms/%d/memberships/", v)
			case "microcosm":
				return fmt.Sprintf("/microcosms/%d/", v)
			case "microcosm-create":
				return fmt.Sprintf("/microcosms/%d/create/microcosm/", v)
			case "profile":
				return fmt.Sprintf("/profiles/%d/", v)
			case "profile-edit":
				return fmt.Sprintf("/profiles/%d/edit/", v)
			}
		default:
		}

	}

	return ""
}
