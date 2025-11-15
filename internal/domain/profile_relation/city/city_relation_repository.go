package city

import (
	"context"
)

type CityRelationRepository interface {
	CreateCityRelation(ctx context.Context, city *CityRelation) (*CityRelation, error)
}
