package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(
	router *gin.Engine,
	roleHandler *handler.RoleHandler,
	authMiddleware gin.HandlerFunc,
) {
	roles := router.Group("/roles")
	roles.Use(authMiddleware)

	{
		roles.GET("", roleHandler.FindAll)
		roles.GET("/:role_id", roleHandler.FindByID)
	}
}
