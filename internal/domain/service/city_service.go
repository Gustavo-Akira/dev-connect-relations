package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
)

type CityService struct {
	cityRepo repository.CityRepository
}

func NewCityService(cityRepo repository.CityRepository) *CityService {
	return &CityService{cityRepo: cityRepo}
}

func (s *CityService) CreateCity(ctx context.Context, city entities.City) (entities.City, error) {
	return s.cityRepo.CreateCity(ctx, city)
}
