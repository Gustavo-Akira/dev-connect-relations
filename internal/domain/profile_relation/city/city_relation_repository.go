package city

import (
	"context"
	"devconnectrelations/internal/domain/recommendation"
)

type CityRelationRepository interface {
	CreateCityRelation(ctx context.Context, city *CityRelation) (*CityRelation, error)
	JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error)
}
