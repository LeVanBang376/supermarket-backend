package routes

import (
	"supermarket-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(
	router *gin.Engine,
	orderHandler *handler.OrderHandler,
	orderItemHandler *handler.OrderItemHandler,
	paymentHandler *handler.PaymentHandler,
	authMiddleware gin.HandlerFunc,
) {
	orders := router.Group("/orders")
	orders.Use(authMiddleware)
	{
		// Orders
		orders.POST(
			"",
			orderHandler.Create,
		)

		orders.GET(
			"",
			orderHandler.FindAll,
		)

		orders.GET(
			"/:order_id",
			orderHandler.FindByID,
		)

		orders.PATCH(
			"/:order_id/cancel",
			orderHandler.Cancel,
		)

		// Order items
		orders.POST(
			"/:order_id/items",
			orderItemHandler.Create,
		)

		orders.GET(
			"/:order_id/items",
			orderItemHandler.FindByOrderID,
		)

		orders.GET(
			"/:order_id/items/:sku_barcode",
			orderItemHandler.FindByID,
		)

		orders.PATCH(
			"/:order_id/items/:sku_barcode",
			orderItemHandler.Update,
		)

		orders.DELETE(
			"/:order_id/items/:sku_barcode",
			orderItemHandler.Delete,
		)

		// Payments
		orders.POST(
			"/:order_id/payments",
			paymentHandler.Create,
		)

		orders.GET(
			"/:order_id/payments",
			paymentHandler.FindByOrderID,
		)
	}

	// Payments
	payments := router.Group("/payments")
	payments.Use(authMiddleware)
	{
		payments.GET(
			"/:payment_id",
			paymentHandler.FindByID,
		)
	}
}
