package relation_controller

import (
	relation_dto "devconnectrelations/internal/adapters/inbound/rest/relation/dto"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RelationController struct {
	service *service.RelationService
}

func CreateNewRelationsController(service service.RelationService) *RelationController {
	return &RelationController{service: &service}
}

func (c *RelationController) CreateRelation(ctx *gin.Context) {
	var createDTO relation_dto.CreateRelationDTO
	if err := ctx.ShouldBind(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	relation, formatError := entities.NewRelation(createDTO.FromId, createDTO.TargetId, entities.RelationType(strings.ToUpper(createDTO.RelationType)), entities.RelationStatusPending)
	if formatError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": formatError.Error()})
		return
	}
	savedRelation, createError := c.service.CreateRelation(ctx, *relation)
	if createError != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": createError.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"relation": savedRelation})
}
