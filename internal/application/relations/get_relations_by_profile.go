package relations

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/relation"
)

type IGetRelationsPaged interface {
	Execute(
		ctx context.Context,
		in GetRelationsPagedInput,
	) (*GetRelationsPagedOutput, error)
}

type GetRelationsPaged struct {
	Repo relation.RelationsRepository
}

type GetRelationsPagedInput struct {
	FromID int64
	Page   int64
}

type GetRelationsPagedOutput struct {
	Data        []relation.Relation
	Page        int64
	TotalItems  int64
	TotalPages  int64
	HasNext     bool
	HasPrevious bool
}

func (uc *GetRelationsPaged) Execute(
	ctx context.Context,
	in GetRelationsPagedInput,
) (*GetRelationsPagedOutput, error) {

	const limit int64 = 20
	offset := in.Page * limit

	items, err := uc.Repo.GetAllRelationsByFromId(ctx, in.FromID, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := uc.Repo.CountRelationsByFromId(ctx, in.FromID)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &GetRelationsPagedOutput{
		Data:        items,
		Page:        in.Page,
		TotalItems:  total,
		TotalPages:  totalPages,
		HasNext:     in.Page+1 < totalPages,
		HasPrevious: in.Page > 0,
	}, nil
}
