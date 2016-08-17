// Portions Copyright 2016 The Hugo Authors

package ui

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"

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
	funcs["relTime"] = relTime
	funcs["reverseLinks"] = reverseLinks
	funcs["stat"] = stat
	funcs["title"] = strings.Title
	funcs["trunc"] = trunc
	funcs["url"] = urlBuilder
	funcs["utcRFC"] = utcRFC

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

func intcomma(value interface{}) (string, error) {
	switch v := value.(type) {
	case float32:
		return humanize.Commaf(float64(v)), nil
	case float64:
		return humanize.Commaf(v), nil
	case int:
		return humanize.Comma(int64(v)), nil
	case int32:
		return humanize.Comma(int64(v)), nil
	case int64:
		return humanize.Comma(v), nil
	default:
		return "", fmt.Errorf("value was not a number")
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

func relTime(d time.Time) string {
	return humanize.Time(d)
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

func trunc(s string, length int) string {
	if s == "" {
		return s
	}

	if len(s) < length {
		return s
	}

	return s[0:length] + "..."
}

// TODO: this is dangerous, no checking of args length or types
func urlBuilder(urlKey string, args ...interface{}) (string, error) {
	switch urlKey {
	case "conversation-create":
		return fmt.Sprintf("/microcosms/%d/create/conversation/", args[0]), nil
	case "event-create":
		return fmt.Sprintf("/microcosms/%d/create/event/", args[0]), nil
	case "home":
		return "/", nil
	case "huddle-list":
		return "/huddles/", nil
	case "legal-list":
		return "/legal/", nil
	case "legal":
		return fmt.Sprintf("/legal/%s", args[0]), nil
	case "login":
		return "/login/", nil
	case "logout":
		return "/logout/", nil
	case "memberships-list":
		return fmt.Sprintf("/microcosms/%d/memberships/", args[0]), nil
	case "microcosm":
		return fmt.Sprintf("/microcosms/%d/", args[0]), nil
	case "microcosm-create":
		return fmt.Sprintf("/microcosms/%d/create/microcosm/", args[0]), nil
	case "profile":
		return fmt.Sprintf("/profiles/%d", args[0]), nil
	case "profile-edit":
		return fmt.Sprintf("/profiles/%d/edit/", args[0]), nil
	case "profile-list":
		return "/profiles/", nil
	case "today":
		return "/today/", nil
	case "update-list":
		return "/updates/", nil
	case "update-settings":
		return "/updates/settings/", nil
	}

	return "", fmt.Errorf("no URL found for '%s'", urlKey)
}

func utcRFC(d time.Time) string {
	return d.UTC().Format(time.RFC3339)
}
