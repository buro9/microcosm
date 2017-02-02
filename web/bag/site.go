package bag

import (
	"context"

	"github.com/buro9/microcosm/models"
)

func SetSite(ctx context.Context, site *models.Site) context.Context {
	return context.WithValue(ctx, contextSite, site)
}

func GetSite(ctx context.Context) *models.Site {
	site, _ := ctx.Value(contextSite).(*models.Site)
	return site
}
