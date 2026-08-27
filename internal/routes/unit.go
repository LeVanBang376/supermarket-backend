package routes

import (
	"supermarket-backend/internal/handler"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterUnitRoutes(
	router *gin.Engine,
	unitHandler *handler.UnitHandler,
	authMiddleware gin.HandlerFunc,
) {
	units := router.Group("/units")
	units.Use(authMiddleware)
	{
		units.POST(
			"",
			middleware.RequireRole(model.RoleAdmin),
			unitHandler.Create,
		)

		units.GET("", unitHandler.FindAll)
		units.GET("/:unit_id", unitHandler.FindByID)

		units.PUT(
			"/:unit_id",
			middleware.RequireRole(model.RoleAdmin),
			unitHandler.Update,
		)

		units.DELETE(
			"/:unit_id",
			middleware.RequireRole(model.RoleAdmin),
			unitHandler.Delete,
		)
	}
}
