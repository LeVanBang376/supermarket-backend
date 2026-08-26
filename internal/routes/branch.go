package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterBranchRoutes(
	router *gin.Engine,
	branchHandler *handler.BranchHandler,
	authMiddleware gin.HandlerFunc,
) {
	branches := router.Group("/branches")
	branches.Use(authMiddleware)
	{
		branches.POST("", branchHandler.Create)
		branches.GET("", branchHandler.FindAll)
		branches.GET("/:branch_id", branchHandler.FindByID)
		branches.PUT("/:branch_id", branchHandler.Update)
		branches.DELETE("/:branch_id", branchHandler.Delete)
	}
}
