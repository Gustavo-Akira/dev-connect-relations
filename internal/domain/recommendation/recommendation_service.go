package recommendation

import (
	"context"
)

type IRecommendationService interface {
	GetRecommendationByProfileId(ctx context.Context, profileID int64) ([]Recommendation, error)
}

type RecommendationService struct {
	RecommendationAlgorithm RecommendationAlgorithm
}

func (recommendation_service *RecommendationService) GetRecommendationByProfileId(ctx context.Context, profileID int64) ([]Recommendation, error) {
	weights := []float64{0.5, 0.3, 0.2}
	recommendations, err := recommendation_service.RecommendationAlgorithm.Run(ctx, weights, profileID)
	if err != nil {
		return nil, err
	}
	return recommendations, nil
}
