// Portions Copyright 2016 The Hugo Authors

package ui

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"

	reflections "gopkg.in/oleiade/reflections.v1"

	humanize "github.com/dustin/go-humanize"
	"github.com/microcosm-cc/microcosm-ui/tpl"
)

func funcMap() template.FuncMap {
	// Inherit all of the Hugo funcs
	funcs := tpl.FuncMap

	// Add our own (or redefine the Hugo ones)
	funcs["exists"] = exists
	funcs["guiURL"] = guiURL
	funcs["hasField"] = hasField
	funcs["intcomma"] = intcomma
	funcs["link"] = link
	funcs["lower"] = strings.ToLower
	funcs["reverseLinks"] = reverseLinks
	funcs["stat"] = stat
	funcs["title"] = strings.Title
	funcs["url"] = urlBuilder

	return funcs
}

// exists determines whether a value is not nil. exists = true if the value
// passed in does not result in nil
func exists(v interface{}) bool {
	if v == nil {
		return false
	}

	g := reflect.ValueOf(v)
	if !g.IsValid() {
		return false
	}
	switch g.Kind() {
	case reflect.String:
		// ok
	default:
		if g.IsNil() {
			return false
		}
	}
	return true
}

func guiURL(u string) string {
	u = strings.Replace(u, `/api/v1`, ``, 1)
	if !strings.HasSuffix(u, `/`) {
		u += `/`
	}
	return u
}

func hasField(s interface{}, fieldName string) bool {
	has, _ := reflections.HasField(s, fieldName)
	return has
}

func intcomma(value interface{}) string {
	switch v := value.(type) {
	case float32:
		return humanize.Commaf(float64(v))
	case float64:
		return humanize.Commaf(v)
	case int:
		return humanize.Comma(int64(v))
	case int32:
		return humanize.Comma(int64(v))
	case int64:
		return humanize.Comma(v)
	default:
		return ""
	}
}

// link returns the link of the given rel
func link(links []Link, rel string) *Link {
	for _, link := range links {
		if link.Rel == rel {
			return &link
		}
	}

	return nil
}

// reverseLinks will reverse a slice of links
func reverseLinks(links []Link) []Link {
	var reversed []Link
	for i := len(links) - 1; i >= 0; i-- {
		reversed = append(reversed, links[i])
	}
	return reversed
}

func stat(stats []Stat, name string) int64 {
	for _, stat := range stats {
		if stat.Metric == name {
			return stat.Value
		}
	}

	return 0
}

// TODO: this is dangerous, no checking of args length or types
func urlBuilder(urlKey string, args ...interface{}) string {
	switch urlKey {
	case "conversation-create":
		return fmt.Sprintf("/microcosms/%d/create/conversation/", args[0])
	case "event-create":
		return fmt.Sprintf("/microcosms/%d/create/event/", args[0])
	case "home":
		return "/"
	case "huddle-list":
		return "/huddles/"
	case "legal-list":
		return "/legal/"
	case "legal":
		return fmt.Sprintf("/legal/%s", args[0])
	case "login":
		return "/login/"
	case "logout":
		return "/logout/"
	case "memberships-list":
		return fmt.Sprintf("/microcosms/%d/memberships/", args[0])
	case "microcosm-create":
		return fmt.Sprintf("/microcosms/%d/create/microcosm/", args[0])
	case "profile":
		return fmt.Sprintf("/profiles/%d", args[0])
	case "profile-edit":
		return fmt.Sprintf("/profiles/%d/edit/", args[0])
	case "profile-list":
		return "/profiles/"
	case "today":
		return "/today/"
	case "update-list":
		return "/updates/"
	case "update-settings":
		return "/updates/settings/"
	}

	return ""
}
