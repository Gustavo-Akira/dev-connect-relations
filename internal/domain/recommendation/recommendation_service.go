package recommendation

import (
	"context"
)

type IRecommendationService interface {
	GetRecommendationByProfileId(ctx context.Context, profileID int64) ([]RecommendationReadModel, error)
}

type RecommendationService struct {
	RecommendationAlgorithm RecommendationAlgorithm
	Read                    ReadModelRepository
}

func (recommendation_service *RecommendationService) GetRecommendationByProfileId(ctx context.Context, profileID int64) ([]RecommendationReadModel, error) {
	weights := []float64{0.5, 0.3, 0.2}
	scores, err := recommendation_service.RecommendationAlgorithm.Run(ctx, weights, profileID)
	if err != nil {
		return nil, err
	}
	return recommendation_service.Read.EnrichRecommendations(ctx, scores)
}
