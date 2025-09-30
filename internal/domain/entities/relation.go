package entities

import "errors"

type RelationType string

const (
	RelationFriend RelationType = "FRIEND"
	RelationFollow RelationType = "FOLLOW"
	RelationBlock  RelationType = "BLOCK"
)

type Relation struct {
	FromID int32
	ToID   int32
	Type   RelationType
}

func NewRelation(fromID int32, toID int32, relationType RelationType) (*Relation, error) {
	if fromID == toID {
		return nil, errors.New("cannot create relation with self")
	}

	if relationType == "" {
		return nil, errors.New("relation type is required")
	}

	return &Relation{
		FromID: fromID,
		ToID:   toID,
		Type:   relationType,
	}, nil
}
