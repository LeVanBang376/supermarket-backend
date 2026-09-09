package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/order_item"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.OrderItemResponse{}

type OrderItemHandler struct {
	service *order_item.Service
}

func NewOrderItemHandler(
	service *order_item.Service,
) *OrderItemHandler {
	return &OrderItemHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Add item to order
// @Description  Add a SKU to an order
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Param        request body dto.AddOrderItemRequest true "Add order item"
// @Success      201 {object} dto.OrderItemResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/items [post]
func (h *OrderItemHandler) Create(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	var req dto.AddOrderItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	actorID := middleware.GetUserID(c)

	if actorID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	item, err := h.service.Create(
		c.Request.Context(),
		actorID,
		orderID,
		&req,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Order or SKU not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Order cannot be modified in its current status",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusCreated,
		"Add item to order successfully",
		item,
	)
}

// FindByID godoc
// @Summary      Get order item
// @Description  Get an order item by order ID and SKU barcode
// @Tags         order-items
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Param        sku_barcode path string true "SKU barcode"
// @Success      200 {object} dto.OrderItemResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/items/{sku_barcode} [get]
func (h *OrderItemHandler) FindByID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	skuBarcode := c.Param("sku_barcode")

	item, err := h.service.FindByID(
		c.Request.Context(),
		orderID,
		skuBarcode,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Order item not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get order item successfully",
		item,
	)
}

// FindByOrderID godoc
// @Summary      Get order items
// @Description  Get all items of an order
// @Tags         order-items
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Success      200 {array} dto.OrderItemResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/items [get]
func (h *OrderItemHandler) FindByOrderID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	items, err := h.service.FindByOrderID(
		c.Request.Context(),
		orderID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Order not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get order items successfully",
		items,
	)
}

// Update godoc
// @Summary      Update order item
// @Description  Update the quantity of an order item
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Param        sku_barcode path string true "SKU barcode"
// @Param        request body dto.UpdateOrderItemRequest true "Update order item"
// @Success      200 {object} dto.OrderItemResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/items/{sku_barcode} [patch]
func (h *OrderItemHandler) Update(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	skuBarcode := c.Param("sku_barcode")

	var req dto.UpdateOrderItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	actorID := middleware.GetUserID(c)

	if actorID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	item, err := h.service.Update(
		c.Request.Context(),
		actorID,
		orderID,
		skuBarcode,
		&req,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Order item not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Order cannot be modified in its current status",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Update order item successfully",
		item,
	)
}

// Delete godoc
// @Summary      Delete order item
// @Description  Delete an item from an order
// @Tags         order-items
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Param        sku_barcode path string true "SKU barcode"
// @Success      200
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/items/{sku_barcode} [delete]
func (h *OrderItemHandler) Delete(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	skuBarcode := c.Param("sku_barcode")

	actorID := middleware.GetUserID(c)

	if actorID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		actorID,
		orderID,
		skuBarcode,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Order item not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Order cannot be modified in its current status",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Delete order item successfully",
	)
}
