package rest

import (
	profile_dto "devconnectrelations/internal/adapters/inbound/rest/profile/dto"
	domain "devconnectrelations/internal/domain/profile"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	service *domain.ProfileService
}

func CreateNewProfileController(service domain.ProfileService) *ProfileController {
	return &ProfileController{service: &service}
}

func (c *ProfileController) CreateProfile(ctx *gin.Context) {
	var createDTO profile_dto.CreateProfileDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile, err := domain.NewProfile(createDTO.Id, createDTO.Name)
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

func (c *ProfileController) DeleteProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedInt, parsedError := strconv.ParseInt(id, 10, 64)
	if parsedError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedError.Error()})
		return
	}
	err := c.service.DeleteProfile(ctx, int64(parsedInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (c *ProfileController) GetProfileByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedInt, parsedError := strconv.ParseInt(id, 10, 64)
	if parsedError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedError.Error()})
		return
	}
	profile, err := c.service.GetProfileByID(ctx, int64(parsedInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"profile": profile})
}
