package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type CityRepository interface {
	CreateCity(ctx context.Context, city entities.City) (entities.City, error)
	GetCityByFullName(ctx context.Context, fullName string) (*entities.City, error)
}
