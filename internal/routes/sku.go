package routes

import (
	"supermarket-backend/internal/handler"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterSKURoutes(
	router *gin.Engine,
	skuHandler *handler.SKUHandler,
	authMiddleware gin.HandlerFunc,
) {
	skus := router.Group("/skus")
	skus.Use(authMiddleware)
	{
		skus.POST(
			"",
			middleware.RequireRole(model.RoleAdmin),
			skuHandler.Create,
		)

		skus.GET("", skuHandler.FindAll)
		skus.GET("/:sku_barcode", skuHandler.FindByID)

		skus.PUT(
			"/:sku_barcode",
			middleware.RequireRole(model.RoleAdmin),
			skuHandler.Update,
		)

		skus.DELETE(
			"/:sku_barcode",
			middleware.RequireRole(model.RoleAdmin),
			skuHandler.Delete,
		)
	}
}
