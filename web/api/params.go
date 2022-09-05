package api

import (
	"context"
	"net/url"
)

type Params struct {
	Ctx        context.Context
	PathPrefix string
	ID         int64
	PathSuffix string
	Q          url.Values
}
