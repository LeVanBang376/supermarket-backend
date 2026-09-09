package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterWebsocketRoutes(
	router *gin.Engine,
	websocketHandler *handler.WebsocketHandler,
	authMiddleware gin.HandlerFunc,
) {
	ws := router.Group("/ws")
	ws.Use(authMiddleware)
	{
		ws.GET("", websocketHandler.Handle)
	}
}
