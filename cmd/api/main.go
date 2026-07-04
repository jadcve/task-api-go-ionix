package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/config"
	"task-api-go-ionix/internal/database"
	"task-api-go-ionix/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbPool, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	router := gin.Default()

	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.GetHealth)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
