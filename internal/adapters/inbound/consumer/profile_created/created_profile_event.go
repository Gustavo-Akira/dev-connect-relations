package profile_created

type CreatedProfileEvent struct {
	Id      int64
	Name    string
	Stack   []string
	City    string
	Country string
	State   string
}
