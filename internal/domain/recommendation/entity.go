package recommendation

type Recommendation struct {
	ID    int64
	Score float64
}

func CreateRecommendation(id int64, score float64) *Recommendation {
	return &Recommendation{
		ID:    id,
		Score: score,
	}
}
