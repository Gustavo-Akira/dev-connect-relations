package stack

import (
	"context"
	"errors"
)

type StackService struct {
	repository StackRepository
}

func CreateStackService(repo StackRepository) *StackService {
	return &StackService{repository: repo}
}

func (s *StackService) CreateStack(ctx context.Context, stackName string) (Stack, error) {
	stack, stackError := NewStack(stackName)
	if stackError != nil {
		return Stack{}, stackError
	}
	old, err := s.repository.GetStackByName(ctx, stack.Name)
	if err != nil {
		return Stack{}, err
	}
	if old != (Stack{}) {
		return Stack{}, errors.New("stack already exists")
	}
	savedStack, saveError := s.repository.CreateStack(ctx, stack)
	return savedStack, saveError
}

func (s *StackService) GetStackByName(ctx context.Context, name string) (Stack, error) {
	stack, err := NewStack(name)
	if err != nil {
		return Stack{}, err
	}
	return s.repository.GetStackByName(ctx, stack.Name)
}

func (s *StackService) DeleteStack(ctx context.Context, name string) error {
	return s.repository.DeleteStack(ctx, name)
}
