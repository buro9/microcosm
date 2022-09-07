package errors

import "fmt"

var (
	// ErrFlagsNoAPIDomain is returned when the API Domain is not defined which
	// means that the API client does not know where to send requests
	ErrFlagsNoAPIDomain = fmt.Errorf(
		`ApiDomain must be provided as -apiDomain or $APIDOMAIN`,
	)

	// ErrNot200 is returned when the API does not return a HTTP 200 when it is
	// expected to. This usually means that there is an issue with the API
	// server or that there is a network issue.
	ErrNot200 = fmt.Errorf(
		`Expected HTTP 200`,
	)

	// ErrAccessTokenExpected is returned when we do not have an access_token to
	// make an authenticated request with.
	ErrAccessTokenExpected = fmt.Errorf(
		`access_token expected`,
	)

	// ErrClientSecretNotConfigured is returned when the OAuth2.0 client secret
	// is not defined and this prevents us from talking to the API server.
	ErrClientSecretNotConfigured = fmt.Errorf(
		`clientSecret not configured`,
	)

	// ErrCodeRequired is returned when we are using Auth0 for login and yet
	// the code that their API defines is not defined.
	ErrCodeRequired = fmt.Errorf(
		`code required for Auth0Login`,
	)

	// ErrSecureCookieMustExist is returned when we are using Gorilla SecureCookie
	// but it has not been initialised.
	ErrSecureCookieMustExist = fmt.Errorf(
		`SecureCookie must exist`,
	)

)
