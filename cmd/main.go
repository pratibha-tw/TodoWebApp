package main

import (
	"todoapp/config"
	"todoapp/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	var cfg = config.Config{}
	config.GetConfigs(&cfg)
	router.RegisterRoutes(engine, cfg)
	engine.Run("localhost:8080")
}
