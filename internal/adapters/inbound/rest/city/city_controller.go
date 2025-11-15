package city_rest

import (
	"devconnectrelations/internal/adapters/inbound/rest/city/dto"
	"devconnectrelations/internal/domain/city"
	"strings"

	"github.com/gin-gonic/gin"
)

type CityController struct {
	cityService city.CityService
}

func CreateNewCityController(cityService city.CityService) *CityController {
	return &CityController{
		cityService: cityService,
	}
}

func (cc *CityController) CreateCity(ctx *gin.Context) {
	var createDTO dto.CreateCityDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	city, err := cc.cityService.CreateCity(ctx.Request.Context(), *city.NewCity(createDTO.Name, createDTO.Country, createDTO.State))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"city": city})
}

func (cc CityController) GetCityByFullName(ctx *gin.Context) {
	fullName := ctx.Param("fullName")
	city, err := cc.cityService.GetCityByFullName(ctx, fullName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, city)
}
