package bag

import (
	"context"

	"github.com/buro9/microcosm/models"
)

func SetProfile(ctx context.Context, profile *models.Profile) context.Context {
	return context.WithValue(ctx, contextProfile, profile)
}

func GetProfile(ctx context.Context) *models.Profile {
	user, _ := ctx.Value(contextProfile).(*models.Profile)
	return user
}
