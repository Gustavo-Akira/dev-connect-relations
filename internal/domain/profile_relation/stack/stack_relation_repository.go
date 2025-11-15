package stack

import (
	"context"
)

type StackRelationRepository interface {
	CreateStackRelation(ctx context.Context, stackRelation *StackRelation) (*StackRelation, error)
	DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error
}
