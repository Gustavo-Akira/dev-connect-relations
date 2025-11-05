package stackrelation

import (
	"context"
	"devconnectrelations/internal/adapters/inbound/rest/stack_relation/dto"
	"devconnectrelations/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type StackRelationController struct {
	service *service.StackRelationService
}

func CreateNewStackRelationController(svc *service.StackRelationService) *StackRelationController {
	return &StackRelationController{service: svc}
}

func (c *StackRelationController) CreateStackRelation(ctx *gin.Context) {
	var request dto.CreateStackRelationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	stackRelation, err := c.service.CreateStackRelation(context.Background(), request.StackName, request.ProfileID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, stackRelation)
}

func (c *StackRelationController) DeleteStackRelation(ctx *gin.Context) {
	var request dto.DeleteStackRelationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.service.DeleteStackRelation(context.Background(), request.StackName, request.ProfileID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Stack relation deleted successfully"})
}
