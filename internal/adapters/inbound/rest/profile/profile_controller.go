package rest

import (
	profile_dto "devconnectrelations/internal/adapters/inbound/rest/profile/dto"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	service *service.ProfileService
}

func CreateNewProfileController(service service.ProfileService) *ProfileController {
	return &ProfileController{service: &service}
}

func (c *ProfileController) CreateProfile(ctx *gin.Context) {
	var createDTO profile_dto.CreateProfileDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := entities.NewProfile(createDTO.Id, createDTO.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context := ctx.Request.Context()
	result, err := c.service.CreateProfile(context, profile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"profile": result})
}
