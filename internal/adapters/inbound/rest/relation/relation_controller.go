package relation_controller

import (
	relation_dto "devconnectrelations/internal/adapters/inbound/rest/relation/dto"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"
	"net/http"
	"strconv"
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

func (c *RelationController) GetAllRelationsByFromId(ctx *gin.Context) {
	fromId := ctx.Param("fromId")
	parsedInt, parsedError := strconv.ParseInt(fromId, 10, 32)
	if parsedError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedError.Error()})
		return
	}
	relations, err := c.service.GetAllRelationsByFromId(ctx, int64(parsedInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"relations": relations})
}

func (c *RelationController) AcceptRelation(ctx *gin.Context) {
	fromId := ctx.Param("fromId")
	toId := ctx.Param("toId")
	parsedFromId, parsedFromError := strconv.ParseInt(fromId, 10, 32)
	parsedToId, parsedToError := strconv.ParseInt(toId, 10, 32)
	if parsedFromError != nil || parsedToError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromId or toId"})
		return
	}
	err := c.service.AcceptRelation(ctx, int64(parsedFromId), int64(parsedToId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Relation accepted successfully"})
}
