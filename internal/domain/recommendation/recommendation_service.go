package recommendation

import (
	"context"
)

type RecommendationService struct {
	RecommendationAlgorithm RecommendationAlgorithm
}

func (recommendation_service *RecommendationService) GetRecommendationByProfileId(ctx context.Context, profileID int64) ([]Recommendation, error) {
	recommendations, err := recommendation_service.RecommendationAlgorithm.Run(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return recommendations, nil
}
