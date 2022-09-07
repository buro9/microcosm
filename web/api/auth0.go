package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/opts"
)

// Auth0Login consumes the auth0 code and will use the microcosm API to fetch
// a microcosm access_token
func Auth0Login(ctx context.Context, code string, state string) (string, int, error) {
	if opts.ClientSecret == nil || *opts.ClientSecret == "" {
		return "", 0, errors.ErrClientSecretNotConfigured
	}

	if code == "" {
		return "", 0, errors.ErrCodeRequired
	}

	type Auth0 struct {
		Code         string `json:"Code"`
		State        string `json:"State"`
		ClientSecret string `json:"ClientSecret"`
	}
	auth0 := Auth0{
		Code:         code,
		State:        state,
		ClientSecret: *opts.ClientSecret,
	}

	resp, err := apiPost(Params{Ctx: ctx, Type: "auth0"}, auth0)
	if err != nil {
		log.Printf("apiPost(`auth0`): %s", err.Error())
		return "", resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("err: not 200")
		return "", resp.StatusCode, errors.ErrNot200
	}

	type Resp struct {
		AccessToken string `json:"data"`
	}
	var r Resp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", resp.StatusCode, err
	}

	return r.AccessToken, resp.StatusCode, nil
}
