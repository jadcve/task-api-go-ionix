package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/response"
	"task-api-go-ionix/internal/security"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			response.Error(c, 401, "Unauthorized", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, 401, "Unauthorized", nil)
			c.Abort()
			return
		}

		claims, err := security.ValidateToken(parts[1], m.jwtSecret)
		if err != nil {
			response.Error(c, 401, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
