package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *entities.Profile) (entities.Profile, error)
	DeleteProfile(ctx context.Context, profileId int64) error
}
