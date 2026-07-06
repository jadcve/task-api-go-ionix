package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"task-api-go-ionix/internal/response"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	response.Success(c, 200, "Service is healthy", gin.H{
		"service": "ok",
	})
}

func (h *HealthHandler) GetDatabaseHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		response.Error(c, 503, "Database connection is unhealthy", gin.H{
			"database": "unavailable",
		})
		return
	}

	response.Success(c, 200, "Database connection is healthy", gin.H{
		"database": "ok",
	})
}
