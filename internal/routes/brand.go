package routes

import (
	"supermarket-backend/internal/handler"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterBrandRoutes(
	router *gin.Engine,
	brandHandler *handler.BrandHandler,
	authMiddleware gin.HandlerFunc,
) {
	brands := router.Group("/brands")
	brands.Use(authMiddleware)
	{
		brands.POST(
			"",
			middleware.RequireRole(model.RoleAdmin),
			brandHandler.Create,
		)

		brands.GET("", brandHandler.FindAll)
		brands.GET("/:brand_id", brandHandler.FindByID)

		brands.PUT(
			"/:brand_id",
			middleware.RequireRole(model.RoleAdmin),
			brandHandler.Update,
		)

		brands.DELETE(
			"/:brand_id",
			middleware.RequireRole(model.RoleAdmin),
			brandHandler.Delete,
		)
	}
}
