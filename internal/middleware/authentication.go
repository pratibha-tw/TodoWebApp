package middleware

import (
	"errors"
	"net/http"
	"strings"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, errors.New("authorization header is missing"))
		ctx.Abort()
		return
	}
	// Token is usually in the format: Bearer <token>
	tokenString := strings.Split(authHeader, " ")[1]
	ctx.Set("token", tokenString)

	claims, err := userservice.VerifyJWT(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.New("user is not authorized"))
		ctx.Abort()
		return
	}
	ctx.Set("user_id", claims.UserId)
	ctx.Next()
}
