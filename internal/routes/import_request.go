package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterImportRequestRoutes(
	router *gin.Engine,
	importRequestHandler *handler.ImportRequestHandler,
	importRequestProductHandler *handler.ImportRequestProductHandler,
	importRequestToteHandler *handler.ImportRequestToteHandler,
	authMiddleware gin.HandlerFunc,
) {
	importRequests := router.Group("/import-requests")
	importRequests.Use(authMiddleware)
	{
		// Import requests
		importRequests.POST(
			"",
			importRequestHandler.Create,
		)

		importRequests.GET(
			"",
			importRequestHandler.FindAll,
		)

		importRequests.GET(
			"/:request_id",
			importRequestHandler.FindByID,
		)

		importRequests.PATCH(
			"/:request_id",
			importRequestHandler.Update,
		)

		importRequests.PATCH(
			"/:request_id/status",
			importRequestHandler.UpdateStatus,
		)

		importRequests.POST(
			"/:request_id/confirm",
			importRequestHandler.ConfirmImport,
		)

		// Products
		importRequests.GET(
			"/:request_id/products",
			importRequestProductHandler.FindByRequestID,
		)

		// Totes
		importRequests.POST(
			"/:request_id/totes/:tote_barcode",
			importRequestToteHandler.Create,
		)

		importRequests.GET(
			"/:request_id/totes",
			importRequestToteHandler.FindAllByRequestID,
		)

		importRequests.GET(
			"/:request_id/totes/:tote_barcode",
			importRequestToteHandler.FindByBarcode,
		)

		importRequests.DELETE(
			"/:request_id/totes/:tote_barcode",
			importRequestToteHandler.Delete,
		)

		// Products of a tote
		importRequests.GET(
			"/:request_id/totes/:tote_barcode/products",
			importRequestProductHandler.FindByRequestAndToteBarcode,
		)
	}
}
