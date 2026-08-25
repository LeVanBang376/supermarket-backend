package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)

		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("/me", authHandler.Me)
		}
	}
}
