package routes

import (
	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/handler"
	"task-api-go-ionix/internal/middleware"
)

func RegisterRoutes(
	router *gin.Engine,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
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

	usersGroup := router.Group("/api/users", authMiddleware.RequireAuth(), middleware.RoleMiddleware(string(enums.UserRoleAdmin)))
	{
		usersGroup.POST("", userHandler.CreateUser)
		usersGroup.GET("", userHandler.GetUsers)
		usersGroup.GET("/:id", userHandler.GetUserByID)
		usersGroup.PUT("/:id", userHandler.UpdateUser)
		usersGroup.DELETE("/:id", userHandler.DeleteUser)
	}

	tasksGroup := router.Group("/api/tasks", authMiddleware.RequireAuth(), middleware.RoleMiddleware(string(enums.UserRoleAdmin)))
	{
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("", taskHandler.GetTasks)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}
}
