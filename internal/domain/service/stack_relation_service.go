package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
	"fmt"
)

type StackRelationService struct {
	repository repository.StackRelationRepository
}

func CreateStackRelationService(repo repository.StackRelationRepository) *StackRelationService {
	return &StackRelationService{repository: repo}
}

func (s *StackRelationService) CreateStackRelation(ctx context.Context, stackName string, profileID int64) (*entities.StackRelation, error) {
	stackRelation, stackRelationError := entities.NewStackRelation(stackName, profileID)
	if stackRelationError != nil {
		return nil, stackRelationError
	}
	fmt.Println("Creating stack relation:", stackRelation)
	savedStackRelation, saveError := s.repository.CreateStackRelation(ctx, stackRelation)
	return savedStackRelation, saveError
}

func (s *StackRelationService) DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error {
	return s.repository.DeleteStackRelation(ctx, stackName, profileID)
}
