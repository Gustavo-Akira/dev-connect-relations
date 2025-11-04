package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
	"errors"
)

type ProfileService struct {
	repository repository.ProfileRepository
}

func CreateNewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repository: repo}
}

func (s *ProfileService) CreateProfile(ctx context.Context, profile *entities.Profile) (entities.Profile, error) {
	return s.repository.CreateProfile(ctx, profile)
}

func (s *ProfileService) DeleteProfile(ctx context.Context, id int64) error {
	if id == int64(0) {
		return errors.New("id to delete cannot be 0")
	}
	return s.repository.DeleteProfile(ctx, id)
}

func (s *ProfileService) GetProfileByID(ctx context.Context, profileId int64) (*entities.Profile, error) {
	return s.repository.GetProfileByID(ctx, profileId)
}
