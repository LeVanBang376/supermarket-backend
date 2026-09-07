package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/order"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.OrderResponse{}

type OrderHandler struct {
	service *order.Service
}

func NewOrderHandler(
	service *order.Service,
) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create order
// @Description  Create a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderRequest true "Create order"
// @Success      201 {object} dto.OrderResponse
// @Failure      400 {object} gin.H
// @Failure      401 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	cashierID := middleware.GetUserID(c)

	if cashierID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	order, err := h.service.Create(
		c.Request.Context(),
		&req,
		cashierID,
	)
	if err != nil {
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
		"Create order successfully",
		order,
	)
}

// FindAll godoc
// @Summary      Get all orders
// @Description  Get all orders
// @Tags         orders
// @Produce      json
// @Param        branch_id query string false "Branch ID"
// @Param        status query string false "Order status"
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.OrderResponse
// @Failure      500 {object} gin.H
// @Router       /orders [get]
func (h *OrderHandler) FindAll(c *gin.Context) {
	var query dto.FindAllOrdersQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	pagination := response.NewPagination(c.Request)

	orders, err := h.service.FindAll(
		c.Request.Context(),
		&query,
		pagination,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.PaginatedJSON(
		c.Writer,
		http.StatusOK,
		"Get orders successfully",
		orders,
		pagination,
	)
}

// FindByID godoc
// @Summary      Get order by ID
// @Description  Get an order by ID
// @Tags         orders
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Success      200 {object} dto.OrderResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id} [get]
func (h *OrderHandler) FindByID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	order, err := h.service.FindByID(
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
		"Get order successfully",
		order,
	)
}

// Cancel godoc
// @Summary      Cancel order
// @Description  Cancel an order
// @Tags         orders
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Success      200
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/cancel [patch]
func (h *OrderHandler) Cancel(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	err = h.service.Cancel(
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

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Order cannot be cancelled in its current status",
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
		"Cancel order successfully",
	)
}
