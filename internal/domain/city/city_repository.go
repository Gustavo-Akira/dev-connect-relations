package city

import (
	"context"
)

type CityRepository interface {
	CreateCity(ctx context.Context, city City) (City, error)
	GetCityByFullName(ctx context.Context, fullName string) (*City, error)
}
