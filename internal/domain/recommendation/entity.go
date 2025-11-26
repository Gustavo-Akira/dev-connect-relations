package recommendation

type Recommendation struct {
	ID       int64
	Name     string
	Score    float64
	CityName string
	Stacks   []string
}

func CreateRecommendation(id int64, score float64, name string) *Recommendation {
	return &Recommendation{
		ID:    id,
		Score: score,
		Name:  name,
	}
}
