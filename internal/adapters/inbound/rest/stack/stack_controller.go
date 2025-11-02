package stack_controller

import (
	"devconnectrelations/internal/adapters/inbound/rest/stack/dto"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type StackController struct {
	service service.StackService
}

func CreateNewStackController(svc service.StackService) *StackController {
	return &StackController{service: svc}
}

func (c *StackController) CreateStack(ctx *gin.Context) {
	var createDTO dto.CreateStackDTO
	if err := ctx.ShouldBind(&createDTO); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	stack, err := c.service.CreateStack(ctx, createDTO.Name)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{"stack": stack})
}

func (c *StackController) GetStackByName(ctx *gin.Context) {
	name := ctx.Param("name")
	stack, err := c.service.GetStackByName(ctx, name)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if stack == (entities.Stack{}) {
		ctx.JSON(404, gin.H{"error": "stack not found"})
		return
	}
	ctx.JSON(200, gin.H{"stack": stack})
}

func (c *StackController) DeleteStack(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.service.DeleteStack(ctx, name)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "stack deleted"})
}
