package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type CityRelationRepository interface {
	CreateCityRelation(ctx context.Context, city *entities.CityRelation) (*entities.CityRelation, error)
}
