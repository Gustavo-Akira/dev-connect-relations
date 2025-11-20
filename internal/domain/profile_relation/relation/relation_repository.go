package relation

import (
	"context"
	"devconnectrelations/internal/domain/recommendation"
)

type RelationsRepository interface {
	CreateRelation(context context.Context, relation Relation) (Relation, error)
	GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]Relation, error)
	AcceptRelation(ctx context.Context, fromId int64, toId int64) error
	GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]Relation, error)
	JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error)
}
