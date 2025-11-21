package algorithms

import "devconnectrelations/internal/domain/recommendation"

type RecommendationAlgorithm interface {
	Run(profileId int64) ([]recommendation.Recommendation, error)
}
