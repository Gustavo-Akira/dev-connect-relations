package relation

import (
	"context"
)

type RelationsRepository interface {
	CreateRelation(context context.Context, relation Relation) (Relation, error)
	GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]Relation, error)
	AcceptRelation(ctx context.Context, fromId int64, toId int64) error
	GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]Relation, error)
}
