package stack

import (
	"context"
	"devconnectrelations/internal/domain/recommendation"
)

type StackRelationRepository interface {
	CreateStackRelation(ctx context.Context, stackRelation *StackRelation) (*StackRelation, error)
	DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error
	JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error)
	GetStackRelationByProfileId(ctx context.Context, profileId int64) ([]StackRelation, error)
}
