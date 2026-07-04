package routes

import (
	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/handler"
	"task-api-go-ionix/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	router.GET("/health", healthHandler.GetHealth)

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)
		authGroup.PATCH("/change-password", authMiddleware.RequireAuth(), authHandler.ChangePassword)
		authGroup.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
	}
}
