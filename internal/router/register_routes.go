package router

import (
	"database/sql"
	"fmt"
	"time"
	redisclient "todoapp/internal/database/redis_client"
	"todoapp/internal/handler/todohandler"
	"todoapp/internal/handler/user"
	"todoapp/internal/middleware"
	todorepo "todoapp/internal/repository/todo_repo"
	"todoapp/internal/repository/user_repo"
	"todoapp/internal/service/todoservice"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

func RegisterRoutes(engine *gin.Engine, dbConnect *sql.DB, redisConn *redis.Client) {

	userRepo := user_repo.NewUserRepository(dbConnect)
	userService := userservice.NewUserService(userRepo)
	redisClt := redisclient.NewRedisClient(redisConn)
	userHandler := user.NewUserHandler(userService, redisClt)

	todoRepo := todorepo.NewTodoRepository(dbConnect)
	todoService := todoservice.NewTodoService(todoRepo)
	todo_Handler := todohandler.NewTodoHandler(todoService)

	c := cron.New()
	c.AddFunc("@every 15m", func() {
		fmt.Println("frequency 15 min ", time.Now())
		todoService.GetTasksNearDueDateButNotCompleted()
	})
	c.Start()

	group := engine.Group("todoapp/api")
	{
		//User apis
		group.POST("/register", userHandler.Register)
		group.POST("/login", userHandler.Login)
		group.POST("/logout", userHandler.Logout)
	}
	todoGroup := engine.Group("todoapp/api")
	todoGroup.Use(middleware.AuthenticateMiddleware(redisConn))
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
