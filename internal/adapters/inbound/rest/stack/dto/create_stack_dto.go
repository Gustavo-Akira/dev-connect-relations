package dto

type CreateStackDTO struct {
	Name string `json:"name" binding:"required"`
}
