package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStockRoutes(
	router *gin.Engine,
	stockHandler *handler.StockHandler,
	authMiddleware gin.HandlerFunc,
) {
	stocks := router.Group("/stocks")
	stocks.Use(authMiddleware)
	{
		stocks.GET("", stockHandler.FindAll)
	}
}
