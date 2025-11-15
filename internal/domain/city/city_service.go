package city

import (
	"context"
)

type CityService struct {
	cityRepo CityRepository
}

func NewCityService(cityRepo CityRepository) *CityService {
	return &CityService{cityRepo: cityRepo}
}

func (s *CityService) CreateCity(ctx context.Context, city City) (City, error) {
	return s.cityRepo.CreateCity(ctx, city)
}

func (s *CityService) GetCityByFullName(ctx context.Context, fullName string) (*City, error) {
	return s.cityRepo.GetCityByFullName(ctx, fullName)
}
