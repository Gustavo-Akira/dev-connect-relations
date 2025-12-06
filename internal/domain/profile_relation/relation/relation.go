package relation

import "errors"

type RelationType string
type RelationStatus string

const (
	RelationFriend RelationType = "FRIEND"
	RelationFollow RelationType = "FOLLOW"
	RelationBlock  RelationType = "BLOCK"
)

const (
	RelationStatusPending  RelationStatus = "PENDING"
	RelationStatusAccepted RelationStatus = "ACCEPTED"
	RelationStatusRejected RelationStatus = "REJECTED"
)

type Relation struct {
	FromID          int64
	FromProfileName string
	ToID            int64
	ToProfileName   string
	Type            RelationType
	Status          RelationStatus
}

func NewRelation(fromID int64, toID int64, relationType RelationType, status RelationStatus) (*Relation, error) {
	if fromID == toID {
		return nil, errors.New("cannot create relation with self")
	}

	if relationType == "" {
		return nil, errors.New("relation type is required")
	}

	if relationType == RelationBlock {
		status = RelationStatusAccepted
	}
	return &Relation{
		FromID: fromID,
		ToID:   toID,
		Type:   relationType,
		Status: status,
	}, nil
}
