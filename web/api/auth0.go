package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/buro9/microcosm/web/opts"
)

// Auth0Login consumes the auth0 code and will use the microcosm API to fetch
// a microcosm access_token
func Auth0Login(ctx context.Context, code string, state string) (string, error) {
	if opts.ClientSecret == nil || *opts.ClientSecret == "" {
		return "", ErrClientSecretNotConfigured
	}

	if code == "" {
		return "", ErrCodeRequired
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

	resp, err := apiPost(Params{Ctx: ctx, Endpoint: "auth0"}, auth0)
	if err != nil {
		log.Printf("apiPost(`auth0`): %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("err: not 200")
		return "", ErrNot200
	}

	type Resp struct {
		AccessToken string `json:"data"`
	}
	var r Resp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	return r.AccessToken, nil
}
