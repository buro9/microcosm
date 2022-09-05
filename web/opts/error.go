package opts

import "fmt"

var (
	// ErrClientSecretRequired is returned when we do not have a client secret
	// for OAuth2.0
	ErrClientSecretRequired = fmt.Errorf(
		"microcosmWebAPIClientSecret is a required arg",
	)
)
