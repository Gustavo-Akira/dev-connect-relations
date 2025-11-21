package algorithms

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/city"
	"devconnectrelations/internal/domain/profile_relation/relation"
	"devconnectrelations/internal/domain/profile_relation/stack"
	"devconnectrelations/internal/domain/recommendation"
	"sort"
)

type JaccardAlgorithm struct {
	CityRelationRepository   city.CityRelationRepository
	RelationsRepository      relation.RelationsRepository
	StacksRelationRepository stack.StackRelationRepository
}

func NewJaccardAlgorithm(
	cityRepo city.CityRelationRepository,
	relationsRepo relation.RelationsRepository,
	stacksRepo stack.StackRelationRepository,
) *JaccardAlgorithm {
	return &JaccardAlgorithm{
		CityRelationRepository:   cityRepo,
		RelationsRepository:      relationsRepo,
		StacksRelationRepository: stacksRepo,
	}
}

func (ja *JaccardAlgorithm) Run(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	city_score, city_error := ja.CityRelationRepository.JaccardIndexByProfileId(ctx, profileId)
	if city_error != nil {
		return nil, city_error
	}
	stack_score, stack_error := ja.StacksRelationRepository.JaccardIndexByProfileId(ctx, profileId)
	if stack_error != nil {
		return nil, stack_error
	}

	relation_score, relation_error := ja.RelationsRepository.JaccardIndexByProfileId(ctx, profileId)
	if relation_error != nil {
		return nil, relation_error
	}

	final := combineScores(city_score, stack_score, relation_score)

	return final, nil
}

func combineScores(scoreSets ...[]recommendation.Recommendation) []recommendation.Recommendation {
	combined := make(map[int64]float64)
	weights := []float64{0.5, 0.3, 0.2}
	for i, set := range scoreSets {
		weight := weights[i]
		for _, s := range set {
			combined[s.ID] += (s.Score * weight)
		}
	}

	result := make([]recommendation.Recommendation, 0, len(combined))
	for id, score := range combined {
		result = append(result, recommendation.Recommendation{
			ID:    id,
			Score: score,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result
}
