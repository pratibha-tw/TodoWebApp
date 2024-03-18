package router

import (
	"database/sql"
	"todoapp/internal/handler/user"
	"todoapp/internal/repository/user_repo"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, dbConnect *sql.DB) {

	userRepo := user_repo.NewUserRepository(dbConnect)
	userService := userservice.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	group := engine.Group("todoapp/api")
	{
		group.POST("/register", userHandler.Register)
		group.POST("/login", userHandler.Login)
	}
}
