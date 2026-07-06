package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/config"
	"task-api-go-ionix/internal/database"
	"task-api-go-ionix/internal/handler"
	"task-api-go-ionix/internal/middleware"
	"task-api-go-ionix/internal/repository"
	"task-api-go-ionix/internal/routes"
	"task-api-go-ionix/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	databaseURL := database.BuildDatabaseURL(cfg)

	if err := database.RunMigrations(databaseURL, cfg.MigrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Printf("migrations executed successfully or no changes pending")

	dbPool, err := database.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := database.SeedInitialAdmin(dbPool); err != nil {
		log.Fatalf("failed to seed initial admin: %v", err)
	}

	router := gin.Default()

	healthHandler := handler.NewHealthHandler(dbPool)
	userRepository := repository.NewUserRepository(dbPool)
	taskRepository := repository.NewTaskRepository(dbPool)
	authService := service.NewAuthService(userRepository, cfg.JWTSecret, cfg.JWTExpirationHours)
	userService := service.NewUserService(userRepository)
	taskService := service.NewTaskService(taskRepository)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	routes.RegisterRoutes(router, healthHandler, authHandler, userHandler, taskHandler, authMiddleware)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
