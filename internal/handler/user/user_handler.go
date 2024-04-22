package user

import (
	"net/http"
	"time"
	user_model "todoapp/internal/database/model/user"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userHandler struct {
	userService userservice.UserService
	redisClient *redis.Client
}

// Login implements UserHandler.
func (u userHandler) Login(ctx *gin.Context) {
	var userLogin user_model.UserCredentials
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if jwtToken, err := u.userService.Login(&userLogin); err != nil {
		switch err.Error() {
		case "please provide valid password":
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		case "error in JWT token generation":
			ctx.JSON(http.StatusInternalServerError, err.Error())
		default:
			ctx.JSON(http.StatusBadRequest, "please provide valid username")
		}
	} else {
		//fmt.Println(userLogin.UserId)
		u.redisClient.Set(ctx, jwtToken, userLogin.UserId, time.Minute*5)
		ctx.JSON(http.StatusOK, jwtToken)
	}
}

func (u userHandler) Register(ctx *gin.Context) {
	var userToRegister user_model.User
	if err := ctx.ShouldBindJSON(&userToRegister); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := u.userService.Register(userToRegister); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusCreated, "Registration is successful")
	}
}

func NewUserHandler(userService userservice.UserService, redisClient *redis.Client) UserHandler {
	return &userHandler{userService, redisClient}
}
