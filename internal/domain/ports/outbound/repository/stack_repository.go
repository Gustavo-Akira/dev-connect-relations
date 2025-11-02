package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type StackRepository interface {
	CreateStack(ctx context.Context, name *entities.Stack) (entities.Stack, error)
	GetStackByName(ctx context.Context, name string) (entities.Stack, error)
	DeleteStack(ctx context.Context, name string) error
}
