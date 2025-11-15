package service

import (
	"context"
	"devconnectrelations/internal/domain/city"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
)

type CityRelationService struct {
	cityRelationRepo repository.CityRelationRepository
	cityService      city.CityService
}

func CreateNewCityRelationService(cityRelationRepo repository.CityRelationRepository, cityService *city.CityService) *CityRelationService {
	return &CityRelationService{
		cityRelationRepo: cityRelationRepo,
		cityService:      *cityService,
	}
}

func (crs *CityRelationService) CreateCityRelation(ctx context.Context, cityRelation *entities.CityRelation) (*entities.CityRelation, error) {
	city, err := crs.cityService.GetCityByFullName(ctx, cityRelation.CityFullName)
	if err != nil {
		return nil, err
	}
	cityRelation.CityFullName = city.GetFullName()
	return crs.cityRelationRepo.CreateCityRelation(ctx, cityRelation)
}
