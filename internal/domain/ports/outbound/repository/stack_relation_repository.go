package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type StackRelationRepository interface {
	CreateStackRelation(ctx context.Context, stackRelation *entities.StackRelation) (*entities.StackRelation, error)
}
