package user

import (
	"fmt"
	"net/http"
	user_model "todoapp/internal/database/model/user"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userHandler struct {
	userService userservice.UserService
}

// Login implements UserHandler.
func (u userHandler) Login(ctx *gin.Context) {
	var userToRegister user_model.User
	if err := ctx.ShouldBindJSON(&userToRegister); err != nil {
		fmt.Println("error in binding json")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if res, err := u.userService.Login(userToRegister); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, res)
	}
}

func (u userHandler) Register(ctx *gin.Context) {
	var userToRegister user_model.User
	if err := ctx.ShouldBindJSON(&userToRegister); err != nil {
		fmt.Println("error in binding json")
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

func NewUserHandler(userService userservice.UserService) UserHandler {
	return &userHandler{userService}
}
