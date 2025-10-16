package relation_dto

type CreateRelationDTO struct {
	FromId       int64
	TargetId     int64
	RelationType string
}
