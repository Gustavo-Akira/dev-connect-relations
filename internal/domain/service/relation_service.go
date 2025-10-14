package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
)

type RelationService struct {
	repository repository.RelationsRepository
}

func CreateRelationService(repo repository.RelationsRepository) *RelationService {
	return &RelationService{repository: repo}
}

func (s *RelationService) CreateRelation(ctx context.Context, relation entities.Relation) (entities.Relation, error) {
	return s.repository.CreateRelation(ctx, relation)
}
