package consumer

type CreatedProfileEvent struct {
	Id    int64
	Name  string
	Stack []string
	City  string
}
