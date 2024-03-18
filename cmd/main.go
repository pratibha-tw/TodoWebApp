package main

import (
	"todoapp/config"
	"todoapp/internal/database"
	"todoapp/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	var config = config.Config{}
	dbConnect := database.CreateConnection(config)
	router.RegisterRoutes(engine, dbConnect)
	engine.Run("localhost:8080")
}
