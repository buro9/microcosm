package funcs

import (
	"fmt"
	"strings"

	"github.com/buro9/microcosm/web/opts"
)

// avatarURL returns a URL to an avatar, which is either a Gravatar URL, an
// uploaded custom avatar or a locally hosted placeholder
func avatarURL(u string, subdomain string) string {
	if strings.Contains(u, `gravatar`) {
		return u
	}
	if strings.Contains(u, `files`) {
		return fmt.Sprintf(
			`https://%s.%s%s`,
			subdomain,
			*opts.APIDomain,
			u,
		)
	}
	return `/static/img/avatar.gif`
}
