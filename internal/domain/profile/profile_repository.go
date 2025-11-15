package profile

import (
	"context"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *Profile) (Profile, error)
	DeleteProfile(ctx context.Context, profileId int64) error
	GetProfileByID(ctx context.Context, profileId int64) (*Profile, error)
}
