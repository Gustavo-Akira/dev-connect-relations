package recommendation

type Recommendation struct {
	ID    string
	Score float64
}

func CreateRecommendation(id string, score float64) *Recommendation {
	return &Recommendation{
		ID:    id,
		Score: score,
	}
}
