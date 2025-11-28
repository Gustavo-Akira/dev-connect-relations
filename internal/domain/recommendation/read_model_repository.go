package recommendation

import (
	"context"
)

type ReadModelRepository interface {
	EnrichRecommendations(ctx context.Context, agg []AggregatedScore) ([]RecommendationReadModel, error)
}
