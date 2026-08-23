package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}
}
