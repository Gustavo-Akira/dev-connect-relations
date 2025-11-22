package auth

type AuthClient interface {
	GetProfile(token string) (*int64, error)
}
