package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type StackRepository interface {
	CreateStack(ctx context.Context, name *entities.Stack) (entities.Stack, error)
}
