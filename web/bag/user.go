package bag

import (
	"context"

	"github.com/buro9/microcosm/models"
)

// SetProfile put the current profile associated with the client request into
// the context
func SetProfile(ctx context.Context, profile *models.Profile) context.Context {
	return context.WithValue(ctx, contextProfile, profile)
}

// GetProfile fetches the current profiles associated with the client request
// from the context
func GetProfile(ctx context.Context) *models.Profile {
	user, _ := ctx.Value(contextProfile).(*models.Profile)
	return user
}
