package stack

import (
	"context"
	"fmt"
)

type StackRelationService struct {
	repository StackRelationRepository
}

func CreateStackRelationService(repo StackRelationRepository) *StackRelationService {
	return &StackRelationService{repository: repo}
}

func (s *StackRelationService) CreateStackRelation(ctx context.Context, stackName string, profileID int64) (*StackRelation, error) {
	stackRelation, stackRelationError := NewStackRelation(stackName, profileID)
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
