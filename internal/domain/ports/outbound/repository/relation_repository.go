package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type RelationsRepository interface {
	CreateRelation(context context.Context, relation entities.Relation) (entities.Relation, error)
	GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]entities.Relation, error)
	AcceptRelation(ctx context.Context, fromId int64, toId int64) error
}
