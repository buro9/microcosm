package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	in       string
	expected string
}

func TestApi2ui(t *testing.T) {
	assert.Equal(t, `/microcosms/`, Api2ui(`/api/v1/microcosms`))
	assert.Equal(t, `/conversations/`, Api2ui(`/api/v1/conversations`))
	assert.Equal(t, `/conversations/12345/`, Api2ui(`/api/v1/conversations/12345`))
	assert.Equal(t, `https://example.microco.sm/microcosms/`, Api2ui(`https://example.microco.sm/api/v1/microcosms`))
	assert.Equal(t, `/conversations/?offset=50`, Api2ui(`/api/v1/conversations?offset=50`))
	assert.Equal(t, `/`, Api2ui(``))

	// The only way to trigger a url.Parse error is via a malformed fragment, in
	// this case an incomplete escape
	assert.Equal(t, ``, Api2ui(`#%2`))
}

func TestURL(t *testing.T) {
	assert.Equal(t, `/`, Url(`home`))
	assert.Equal(t, `/huddles/create/`, Url(`huddle-create`))
	assert.Equal(t, `/huddles/`, Url(`huddle-list`))
	assert.Equal(t, `/legal/`, Url(`legal-list`))
	assert.Equal(t, `/login/`, Url(`login`))
	assert.Equal(t, `/logout/`, Url(`logout`))
	assert.Equal(t, `/profiles/`, Url(`profile-list`))
	assert.Equal(t, `/search/`, Url(`search`))
	assert.Equal(t, `/today/`, Url(`today`))
	assert.Equal(t, `/updates/`, Url(`update-list`))
	assert.Equal(t, `/updates/settings/`, Url(`update-settings`))

	assert.Equal(t, `/comments/123/incontext/`, Url(`comment-incontext`, 123))
	assert.Equal(t, `/legal/privacy/`, Url(`legal`, `privacy`))
	assert.Equal(t, `/microcosms/123/`, Url(`microcosm`, 123))
	assert.Equal(t, `/microcosms/123/create/conversation/`, Url(`conversation-create`, 123))
	assert.Equal(t, `/microcosms/123/create/event/`, Url(`event-create`, 123))
	assert.Equal(t, `/microcosms/123/create/microcosm/`, Url(`microcosm-create`, 123))
	assert.Equal(t, `/microcosms/123/memberships/`, Url(`memberships-list`, 123))
	assert.Equal(t, `/profiles/123/`, Url(`profile`, 123))
	assert.Equal(t, `/profiles/123/edit/`, Url(`profile-edit`, 123))

	assert.Equal(t, ``, Url(`microcosm`, `cannot_be_a_string`))
	assert.Equal(t, ``, Url(`this_does_not_exist`))
}
