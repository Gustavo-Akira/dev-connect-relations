package recommendation

type RecommendationReadModel struct {
	ID     int64
	Name   string
	Score  float64
	City   string
	Stacks []string
}

type AggregatedScore struct {
	ID    int64
	Score float64
	Name  string
}
