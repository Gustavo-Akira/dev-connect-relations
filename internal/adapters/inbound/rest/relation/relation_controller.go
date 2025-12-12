package relation_controller

import (
	relation_dto "devconnectrelations/internal/adapters/inbound/rest/relation/dto"
	usecases "devconnectrelations/internal/application/relations"
	"devconnectrelations/internal/domain/profile_relation/relation"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RelationController struct {
	service             relation.IRelationService
	getRelationsUseCase usecases.IGetRelationsPaged
}

func CreateNewRelationsController(service relation.IRelationService, getRelationUseCase usecases.IGetRelationsPaged) *RelationController {
	return &RelationController{service: service, getRelationsUseCase: getRelationUseCase}
}

func (c *RelationController) CreateRelation(ctx *gin.Context) {
	var createDTO relation_dto.CreateRelationDTO

	if err := ctx.ShouldBind(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authError := CompareAndGetUserId(ctx, createDTO.FromId)
	if authError != nil {
		if authError.Error() == "Unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": authError.Error()})
		}
		return
	}

	relation, formatError := relation.NewRelation(createDTO.FromId, createDTO.TargetId, relation.RelationType(strings.ToUpper(createDTO.RelationType)), relation.RelationStatusPending)
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
	page := ctx.DefaultQuery("page", "0")

	parsedPage, parsedPageError := strconv.ParseInt(page, 10, 64)
	if parsedPageError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedPageError.Error()})
		return
	}
	parsedInt, parsedError := strconv.ParseInt(fromId, 10, 64)
	if parsedError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedError.Error()})
		return
	}
	authError := CompareAndGetUserId(ctx, parsedInt)
	if authError != nil {
		if authError.Error() == "Unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"error": authError.Error()})
		}
		return
	}
	relations, err := c.getRelationsUseCase.Execute(ctx, usecases.GetRelationsPagedInput{FromID: int64(parsedInt), Page: parsedPage})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, relations)
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
	authError := CompareAndGetUserId(ctx, parsedFromId)
	if authError != nil {
		if authError.Error() == "Unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": authError.Error()})
		}
		return
	}

	err := c.service.AcceptRelation(ctx, int64(parsedFromId), int64(parsedToId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Relation accepted successfully"})
}

func (c *RelationController) GetAllRelationPendingByFromId(ctx *gin.Context) {
	fromId := ctx.Param("fromId")

	parsedInt, parsedError := strconv.ParseInt(fromId, 10, 64)
	if parsedError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": parsedError.Error()})
		return
	}
	authError := CompareAndGetUserId(ctx, parsedInt)
	if authError != nil {
		if authError.Error() == "Unauthorized" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": authError.Error()})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": authError.Error()})
		}
		return
	}
	relations, err := c.service.GetAllRelationPendingByFromId(ctx, int64(parsedInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"relations": relations})
}

func CompareAndGetUserId(ctx *gin.Context, comparedId int64) error {
	userIDv, exists := ctx.Get("userId")

	if !exists {
		return errors.New("Unauthorized")
	}
	userId := *userIDv.(*int64)
	if userId != comparedId {
		return errors.New("Forbidden " + strconv.FormatInt(userId, 10) + " " + strconv.FormatInt(comparedId, 10))
	}

	return nil
}
