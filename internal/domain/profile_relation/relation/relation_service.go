package relation

import (
	"context"
)

type IRelationService interface {
	CreateRelation(ctx context.Context, relation Relation) (Relation, error)
	GetAllRelationsByFromId(ctx context.Context, fromId int64, page int64) ([]Relation, error)
	AcceptRelation(ctx context.Context, fromId int64, toId int64) error
	GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]Relation, error)
}

type RelationService struct {
	repository RelationsRepository
}

func CreateRelationService(repo RelationsRepository) *RelationService {
	return &RelationService{repository: repo}
}

func (s *RelationService) CreateRelation(ctx context.Context, relation Relation) (Relation, error) {
	return s.repository.CreateRelation(ctx, relation)
}

func (s *RelationService) GetAllRelationsByFromId(ctx context.Context, fromId int64, page int64) ([]Relation, error) {
	var limit int64 = 20
	offset := limit * page
	return s.repository.GetAllRelationsByFromId(ctx, fromId, offset, limit)
}

func (s *RelationService) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	return s.repository.AcceptRelation(ctx, fromId, toId)
}

func (s *RelationService) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]Relation, error) {
	return s.repository.GetAllRelationPendingByFromId(ctx, fromId)
}
