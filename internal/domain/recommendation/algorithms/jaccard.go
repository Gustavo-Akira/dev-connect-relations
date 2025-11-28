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

type NameAndScore struct {
	Name  string
	Score float64
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

func (ja *JaccardAlgorithm) Run(ctx context.Context, weights []float64, profileId int64) ([]recommendation.AggregatedScore, error) {
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

	final := combineScores(weights, city_score, stack_score, relation_score)
	return final, nil
}

func combineScores(weights []float64, scoreSets ...[]recommendation.Recommendation) []recommendation.AggregatedScore {
	combined := make(map[int64]NameAndScore)

	for i, set := range scoreSets {
		weight := weights[i]
		for _, s := range set {
			combined[s.ID] = NameAndScore{
				Score: combined[s.ID].Score + (s.Score * weight),
				Name:  s.Name,
			}
		}
	}

	result := make([]recommendation.AggregatedScore, 0, len(combined))
	for id, score := range combined {
		result = append(result, recommendation.AggregatedScore{
			ID:    id,
			Score: score.Score,
			Name:  score.Name,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result
}
