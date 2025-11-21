package recommendation

import (
	"context"
)

type RecommendationAlgorithm interface {
	Run(ctx context.Context, weights []float64, profileId int64) ([]Recommendation, error)
}
