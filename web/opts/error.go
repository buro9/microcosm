package opts

import "fmt"

var (
	ErrFilesPathNotReadable = fmt.Errorf(
		"microcosmWebFiles must exist and be readable",
	)

	ErrClientSecretRequired = fmt.Errorf(
		"microcosmWebAPIClientSecret is a required arg",
	)
)
