package middlewares

import (
	"net/http"
	"strings"

	authclient "devconnectrelations/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	Client authclient.AuthClient
}

func NewAuthMiddleware(c authclient.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{Client: c}
}

func (a *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("jwt")
		if token == "" || err != nil {
			ctx.Next()
			return
		}
		token = strings.TrimPrefix(token, "")
		user_id, err := a.Client.GetProfile(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		ctx.Set("userId", user_id)
		ctx.Next()
	}
}
