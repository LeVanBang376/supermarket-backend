package middleware

import (
	"net/http"
	"strings"

	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/internal/response"

	"github.com/gin-gonic/gin"
)

const ContextUserKey = "user"

func Auth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Missing authorization header",
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid authorization header",
			)
			c.Abort()
			return
		}

		if parts[0] != "Bearer" {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid auth scheme",
			)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid token",
			)
			c.Abort()
			return
		}

		c.Set(ContextUserKey, claims)

		c.Next()
	}
}

func GetClaims(c *gin.Context) *jwt.Claims {
	claims, exists := c.Get(ContextUserKey)

	if !exists {
		return nil
	}

	result, ok := claims.(*jwt.Claims)
	if !ok {
		return nil
	}

	return result
}
