package recommendation

import (
	"context"
)

type RecommendationAlgorithm interface {
	Run(ctx context.Context, profileId int64) ([]Recommendation, error)
}
