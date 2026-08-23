package middleware

import (
	"net/http"

	"supermarket-backend/internal/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(roleID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)

		if claims == nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			c.Abort()
			return
		}

		if claims.RoleID != roleID {
			response.NonDataJSON(
				c.Writer,
				http.StatusForbidden,
				"You don't have permission to access this resource",
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
