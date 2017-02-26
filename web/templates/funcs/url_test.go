package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	in       string
	expected string
}

func TestAPI2UI(t *testing.T) {
	assert.Equal(t, `/microcosms/`, api2ui(`/api/v1/microcosms`))
	assert.Equal(t, `/conversations/`, api2ui(`/api/v1/conversations`))
	assert.Equal(t, `/conversations/12345/`, api2ui(`/api/v1/conversations/12345`))
	assert.Equal(t, `https://example.microco.sm/microcosms/`, api2ui(`https://example.microco.sm/api/v1/microcosms`))
	assert.Equal(t, `/conversations/?offset=50`, api2ui(`/api/v1/conversations?offset=50`))
	assert.Equal(t, `/`, api2ui(``))

	// The only way to trigger a url.Parse error is via a malformed fragment, in
	// this case an incomplete escape
	assert.Equal(t, ``, api2ui(`#%2`))
}

func TestURL(t *testing.T) {
	assert.Equal(t, `/`, url(`home`))
	assert.Equal(t, `/huddles/`, url(`huddle-list`))
	assert.Equal(t, `/legal/`, url(`legal-list`))
	assert.Equal(t, `/login/`, url(`login`))
	assert.Equal(t, `/logout/`, url(`logout`))
	assert.Equal(t, `/profiles/`, url(`profile-list`))
	assert.Equal(t, `/search/`, url(`search`))
	assert.Equal(t, `/today/`, url(`today`))
	assert.Equal(t, `/updates/`, url(`update-list`))
	assert.Equal(t, `/updates/settings/`, url(`update-settings`))

	assert.Equal(t, `/comments/123/incontext/`, url(`comment-incontext`, 123))
	assert.Equal(t, `/legal/privacy/`, url(`legal`, `privacy`))
	assert.Equal(t, `/microcosms/123/`, url(`microcosm`, 123))
	assert.Equal(t, `/microcosms/123/create/conversation/`, url(`conversation-create`, 123))
	assert.Equal(t, `/microcosms/123/create/event/`, url(`event-create`, 123))
	assert.Equal(t, `/microcosms/123/create/microcosm/`, url(`microcosm-create`, 123))
	assert.Equal(t, `/microcosms/123/memberships/`, url(`memberships-list`, 123))
	assert.Equal(t, `/profiles/123/`, url(`profile`, 123))
	assert.Equal(t, `/profiles/123/edit/`, url(`profile-edit`, 123))

	assert.Equal(t, ``, url(`microcosm`, `cannot_be_a_string`))
	assert.Equal(t, ``, url(`this_does_not_exist`))
}
