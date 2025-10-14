package relation_dto

type CreateRelationDTO struct {
	FromId       int32
	TargetId     int32
	RelationType string
}
