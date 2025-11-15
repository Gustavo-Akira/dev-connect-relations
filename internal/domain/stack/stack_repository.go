package stack

import (
	"context"
)

type StackRepository interface {
	CreateStack(ctx context.Context, name *Stack) (Stack, error)
	GetStackByName(ctx context.Context, name string) (Stack, error)
	DeleteStack(ctx context.Context, name string) error
}
