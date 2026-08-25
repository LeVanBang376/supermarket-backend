package middleware

import (
	"net/http"

	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/internal/response"

	"github.com/gin-gonic/gin"
)

const ContextUserKey = "user"

func Auth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")

		if err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			c.Abort()
			return
		}

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid or expired token",
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
