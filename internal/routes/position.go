package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPositionRoutes(
	router *gin.Engine,
	positionHandler *handler.PositionHandler,
	authMiddleware gin.HandlerFunc,
) {
	positions := router.Group("/positions")
	positions.Use(authMiddleware)

	positions.GET("", positionHandler.FindAll)
	positions.GET("/:position_id", positionHandler.FindByID)
}
