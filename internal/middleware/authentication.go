package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

func AuthenticateMiddleware(redisConn *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, errors.New("authorization header is missing"))
			ctx.Abort()
			return
		}
		// Token is usually in the format: Bearer <token>
		tokenString := strings.Split(authHeader, " ")[1]
		ctx.Set("token", tokenString)

		tokenExist, err := redisConn.Exists(ctx, tokenString).Result()
		if err != nil || tokenExist != 1 {
			ctx.JSON(http.StatusUnauthorized, errors.New("user is not authorized"))
			ctx.Abort()
			return
		}
		userId, _ := redisConn.Get(ctx, tokenString).Result()
		fmt.Println("authentication user id", userId)
		fmt.Printf("type of id = %T\n", userId)
		ctx.Set("user_id", userId)
		// Call the next handler
		ctx.Next()
	}
}
