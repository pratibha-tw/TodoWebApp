package main

import (
	"todoapp/config"
	"todoapp/internal/database"
	"todoapp/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	var cfg = config.Config{}
	config.GetConfigs(&cfg)
	dbConnect := database.CreateConnection(cfg)
	redisConn := database.CreateRedisConnection(cfg)
	router.RegisterRoutes(engine, dbConnect, redisConn)
	engine.Run("localhost:8080")
}
