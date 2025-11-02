package repository

import "devconnectrelations/internal/domain/entities"

type StackRelationRepository interface {
	CreateStackRelation(stackRelation *entities.StackRelation) (*entities.StackRelation, error)
}
