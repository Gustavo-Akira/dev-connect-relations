package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
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
