package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/response"
)

type RoleAuthorizer struct{}

func NewRoleMiddleware() *RoleAuthorizer {
	return &RoleAuthorizer{}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return NewRoleMiddleware().RequireRoles(allowedRoles...)
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return RequireRoles(allowedRoles...)
}
func (m *RoleAuthorizer) RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		normalized := strings.ToUpper(strings.TrimSpace(role))
		allowed[string(enums.UserRole(normalized))] = struct{}{}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			response.Error(c, 401, "Unauthorized", nil)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			response.Error(c, 401, "Unauthorized", nil)
			c.Abort()
			return
		}

		normalizedRole := string(enums.UserRole(strings.ToUpper(strings.TrimSpace(role))))
		if _, ok := allowed[normalizedRole]; !ok {
			response.Error(c, 403, "Forbidden", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
