package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"
)

type RelationsRepository interface {
	CreateRelation(context context.Context, relation entities.Relation) (entities.Relation, error)
}
