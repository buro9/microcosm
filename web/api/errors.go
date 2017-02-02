package api

import "fmt"

var (
	ErrFlagsNoAPIDomain = fmt.Errorf(
		"ApiDomain must be provided as -apiDomain or $APIDOMAIN",
	)
)
