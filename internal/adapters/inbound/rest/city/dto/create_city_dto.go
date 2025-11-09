package dto

type CreateCityDTO struct {
	Name    string `json:"name" binding:"required"`
	State   string `json:"state" binding:"required"`
	Country string `json:"country" binding:"required"`
}
