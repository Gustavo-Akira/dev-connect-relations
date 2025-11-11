package cityrelation

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"devconnectrelations/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type CityRelationController struct {
	cityRelationService service.CityRelationService
}

func CreateNewCityRelationController(cityRelationService service.CityRelationService) *CityRelationController {
	return &CityRelationController{
		cityRelationService: cityRelationService,
	}
}

func (crc *CityRelationController) CreateCityRelation(ctx *gin.Context) {
	var cityRelation entities.CityRelation
	if err := ctx.ShouldBindJSON(&cityRelation); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	createdCityRelation, err := crc.cityRelationService.CreateCityRelation(context.Background(), &cityRelation)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, createdCityRelation)
}
