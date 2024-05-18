package router

import (
	"fmt"
	"sync"
	"time"
	"todoapp/config"
	"todoapp/internal/database"
	"todoapp/internal/database/model/todo"
	redisclient "todoapp/internal/database/redis_client"
	"todoapp/internal/handler/todohandler"
	"todoapp/internal/handler/user"
	emailbody "todoapp/internal/helpers/email_body"
	"todoapp/internal/middleware"
	todorepo "todoapp/internal/repository/todo_repo"
	"todoapp/internal/repository/user_repo"
	"todoapp/internal/service/emailservice"
	"todoapp/internal/service/todoservice"
	"todoapp/internal/service/userservice"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func RegisterRoutes(engine *gin.Engine, cfg config.Config) {

	time.Sleep(10 * time.Second)
	dbConnect := database.CreateConnection(cfg)
	redisConn := database.CreateRedisConnection(cfg)
	database.RunMigration(dbConnect)

	emailservice := emailservice.NewEmailService(cfg.Email)
	userRepo := user_repo.NewUserRepository(dbConnect)
	userService := userservice.NewUserService(userRepo)
	redisClt := redisclient.NewRedisClient(redisConn)
	userHandler := user.NewUserHandler(userService, redisClt)

	todoRepo := todorepo.NewTodoRepository(dbConnect)
	todoService := todoservice.NewTodoService(todoRepo)
	todo_Handler := todohandler.NewTodoHandler(todoService)

	c := cron.New()
	c.AddFunc("@every 5m", func() {
		fmt.Println("frequency 5 min ", time.Now())
		tasks := todoService.GetTasksNearDueDateButNotCompleted()
		var uids []any
		for _, t := range tasks {
			uids = append(uids, t.UserId)
		}
		if len(tasks) > 0 {
			email_list := userRepo.GetEmailIds(uids)
			var wg sync.WaitGroup
			for _, t := range tasks {
				wg.Add(1)
				go func(tk todo.Task) {
					defer wg.Done()
					body := emailbody.Generate_task_notification_email_body(email_list[tk.UserId].Username, tk)
					emailservice.SendNotification(email_list[tk.UserId].Email, body)
				}(t)

			}
		}

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
