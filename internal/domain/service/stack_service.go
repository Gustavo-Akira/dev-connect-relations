package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/ports/outbound/repository"
	"errors"
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
	old, err := s.repository.GetStackByName(ctx, stack.Name)
	if err != nil {
		return entities.Stack{}, err
	}
	if old != (entities.Stack{}) {
		return entities.Stack{}, errors.New("stack already exists")
	}
	savedStack, saveError := s.repository.CreateStack(ctx, stack)
	return savedStack, saveError
}

func (s *StackService) GetStackByName(ctx context.Context, name string) (entities.Stack, error) {
	stack, err := entities.NewStack(name)
	if err != nil {
		return entities.Stack{}, err
	}
	return s.repository.GetStackByName(ctx, stack.Name)
}

func (s *StackService) DeleteStack(ctx context.Context, name string) error {
	return s.repository.DeleteStack(ctx, name)
}
