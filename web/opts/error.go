package opts

import "fmt"

var (
	// ErrFilesPathNotReadable is returned when the static files path that holds
	// templates and resources cannot be read
	ErrFilesPathNotReadable = fmt.Errorf(
		"microcosmWebFiles must exist and be readable",
	)

	// ErrClientSecretRequired is returned when we do not have a client secret
	// for OAuth2.0
	ErrClientSecretRequired = fmt.Errorf(
		"microcosmWebAPIClientSecret is a required arg",
	)
)
