package city

import (
	"context"
	"devconnectrelations/internal/domain/city"
)

type CityRelationService struct {
	cityRelationRepo CityRelationRepository
	cityService      city.CityService
}

func CreateNewCityRelationService(cityRelationRepo CityRelationRepository, cityService *city.CityService) *CityRelationService {
	return &CityRelationService{
		cityRelationRepo: cityRelationRepo,
		cityService:      *cityService,
	}
}

func (crs *CityRelationService) CreateCityRelation(ctx context.Context, cityRelation *CityRelation) (*CityRelation, error) {
	city, err := crs.cityService.GetCityByFullName(ctx, cityRelation.CityFullName)
	if err != nil {
		return nil, err
	}
	cityRelation.CityFullName = city.GetFullName()
	return crs.cityRelationRepo.CreateCityRelation(ctx, cityRelation)
}
