package relation

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/city"
	"devconnectrelations/internal/domain/profile_relation/stack"
	"devconnectrelations/internal/domain/recommendation"
)

type Neo4jRecommendationReadRepository struct {
	cityRepo  city.CityRelationRepository
	stackRepo stack.StackRelationRepository
}

func CreateNeo4jRecommendationRepository(cityRepo city.CityRelationRepository, stackRepo stack.StackRelationRepository) *Neo4jRecommendationReadRepository {
	return &Neo4jRecommendationReadRepository{
		cityRepo:  cityRepo,
		stackRepo: stackRepo,
	}
}

func (r *Neo4jRecommendationReadRepository) EnrichRecommendations(ctx context.Context, agg []recommendation.AggregatedScore) ([]recommendation.RecommendationReadModel, error) {
	ids := make([]int64, 0, len(agg))

	for _, a := range agg {
		ids = append(ids, a.ID)
	}

	cities, err := r.cityRepo.GetCityRelatedToProfileIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	cityMap := make(map[int64]string)
	for _, c := range cities {
		cityMap[c.ProfileID] = c.CityFullName
	}

	stacks, err := r.stackRepo.GetStackRelationByProfileIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	stackMap := make(map[int64][]string)
	for _, s := range stacks {
		stackMap[s.ProfileID] = append(stackMap[s.ProfileID], s.StackName)
	}

	final := make([]recommendation.RecommendationReadModel, 0, len(agg))

	for _, a := range agg {
		id := a.ID

		final = append(final, recommendation.RecommendationReadModel{
			ID:     id,
			Score:  a.Score,
			Name:   a.Name,
			City:   cityMap[id],
			Stacks: stackMap[id],
		})
	}

	return final, nil
}
