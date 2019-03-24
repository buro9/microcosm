package bag

import (
	"context"

	"github.com/buro9/microcosm/models"
)

// SetSite puts the current site into the context
func SetSite(ctx context.Context, site *models.Site) context.Context {
	return context.WithValue(ctx, contextSite, site)
}

// GetSite fetches the curreent site from the context
func GetSite(ctx context.Context) *models.Site {
	site, _ := ctx.Value(contextSite).(*models.Site)
	return site
}
