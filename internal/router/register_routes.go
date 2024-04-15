package router

import (
	"database/sql"
	"todoapp/internal/handler/todohandler"
	"todoapp/internal/handler/user"
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

		//task api
		group.POST("/task/add", todo_Handler.AddTask)
		group.PUT("/task/edit", todo_Handler.UpdateTask)
		group.GET("/task/:id", todo_Handler.GetTaskDetails)
		group.GET("/user/:id/tasks", todo_Handler.GetTodoList)
		group.DELETE("/task/delete/:id", todo_Handler.DeleteTask)
		group.GET("user/:id/notifications", todo_Handler.GetNotifications)
	}
}
