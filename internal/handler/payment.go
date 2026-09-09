package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/payment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.PaymentResponse{}

type PaymentHandler struct {
	service *payment.Service
}

func NewPaymentHandler(
	service *payment.Service,
) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create payment
// @Description  Create a payment for an order
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Param        request body dto.CreatePaymentRequest true "Create payment"
// @Success      201 {object} dto.PaymentResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	payment, err := h.service.Create(
		c.Request.Context(),
		orderID,
		&req,
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
				"Payment cannot be created for this order",
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
		"Create payment successfully",
		payment,
	)
}

// FindByID godoc
// @Summary      Get payment by ID
// @Description  Get a payment by ID
// @Tags         payments
// @Produce      json
// @Param        payment_id path string true "Payment ID"
// @Success      200 {object} dto.PaymentResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /payments/{payment_id} [get]
func (h *PaymentHandler) FindByID(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("payment_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid payment ID",
		)
		return
	}

	payment, err := h.service.FindByID(
		c.Request.Context(),
		paymentID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Payment not found",
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
		"Get payment successfully",
		payment,
	)
}

// FindByOrderID godoc
// @Summary      Get order payments
// @Description  Get all payments of an order
// @Tags         payments
// @Produce      json
// @Param        order_id path string true "Order ID"
// @Success      200 {array} dto.PaymentResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /orders/{order_id}/payments [get]
func (h *PaymentHandler) FindByOrderID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid order ID",
		)
		return
	}

	payments, err := h.service.FindByOrderID(
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
		"Get order payments successfully",
		payments,
	)
}
