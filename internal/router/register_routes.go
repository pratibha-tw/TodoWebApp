package router

import (
	"database/sql"
	"todoapp/internal/handler/todohandler"
	"todoapp/internal/handler/user"
	"todoapp/internal/middleware"
	todorepo "todoapp/internal/repository/todo_repo"
	"todoapp/internal/repository/user_repo"
	"todoapp/internal/service/todoservice"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, dbConnect *sql.DB) {

	userRepo := user_repo.NewUserRepository(dbConnect)
	userService := userservice.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	todoRepo := todorepo.NewTodoRepository(dbConnect)
	todoService := todoservice.NewTodoService(todoRepo)
	todo_Handler := todohandler.NewTodoHandler(todoService)
	go todoService.GetTasksNearDueDateButNotCompleted()
	group := engine.Group("todoapp/api")
	{
		//User apis
		group.POST("/register", userHandler.Register)
		group.POST("/login", userHandler.Login)
	}
	todoGroup := engine.Group("todoapp/api")
	todoGroup.Use(middleware.AuthenticationMiddleware)
	{
		//task api
		todoGroup.POST("/task/add", todo_Handler.AddTask)
		todoGroup.PUT("/task/edit", todo_Handler.UpdateTask)
		todoGroup.GET("/task/:id", todo_Handler.GetTaskDetails)
		todoGroup.GET("/user/:id/tasks", todo_Handler.GetTodoList)
		todoGroup.DELETE("/task/delete/:id", todo_Handler.DeleteTask)
		todoGroup.GET("user/:id/notifications", todo_Handler.GetNotifications)
	}

}
