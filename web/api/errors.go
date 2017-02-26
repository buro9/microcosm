package api

import "fmt"

var (
	ErrFlagsNoAPIDomain = fmt.Errorf(
		"ApiDomain must be provided as -apiDomain or $APIDOMAIN",
	)

	ErrNot200 = fmt.Errorf(
		"Expected HTTP 200",
	)

	ErrAccessTokenExpected = fmt.Errorf(
		"access_token expected",
	)

	ErrClientSecretNotConfigured = fmt.Errorf(
		"clientSecret not configured",
	)

	ErrCodeRequired = fmt.Errorf(
		"code required for Auth0Login",
	)
)
