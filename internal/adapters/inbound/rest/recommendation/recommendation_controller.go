package recommendation

import (
	"devconnectrelations/internal/domain/recommendation"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	service recommendation.IRecommendationService
}

func NewRecommendationController(service recommendation.IRecommendationService) *RecommendationController {
	return &RecommendationController{
		service: service,
	}
}

func (rc *RecommendationController) GetRecommendations(ctx *gin.Context) {
	userID := ctx.Param("userId")
	userIDInt64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	userIDv, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	if *userIDv.(*int64) != userIDInt64 {
		ctx.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	recommendations, recommendationError := rc.service.GetRecommendationByProfileId(ctx, int64(userIDInt64))

	if recommendationError != nil {
		ctx.JSON(500, gin.H{"error": recommendationError.Error()})
		return
	}
	ctx.JSON(200, recommendations)
}
