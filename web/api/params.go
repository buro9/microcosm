package api

import (
	"context"
	"net/url"
)

type Params struct {
	Ctx      context.Context
	Endpoint string
	ID       int64
	Q        url.Values
}
