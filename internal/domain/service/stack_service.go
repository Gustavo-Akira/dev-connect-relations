package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
)

type StackService struct {
	repository repository.StackRepository
}

func CreateStackService(repo repository.StackRepository) *StackService {
	return &StackService{repository: repo}
}

func (s *StackService) CreateStack(ctx context.Context, stackName string) (entities.Stack, error) {
	stack, stackError := entities.NewStack(stackName)
	if stackError != nil {
		return entities.Stack{}, stackError
	}
	savedStack, err := s.repository.CreateStack(ctx, stack)
	return savedStack, err
}
