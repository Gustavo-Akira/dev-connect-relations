package dto

type CreateStackRelationRequest struct {
	StackName string `json:"stack_name" binding:"required"`
	ProfileID int64  `json:"profile_id" binding:"required"`
}
