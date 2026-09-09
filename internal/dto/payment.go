package dto

import (
	"time"

	"github.com/google/uuid"

	"supermarket-backend/internal/model"
)

type CreatePaymentRequest struct {
	Method         string  `json:"method" binding:"required,max=20"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	TransactionRef *string `json:"transaction_ref,omitempty,max=100"`
}

type PaymentResponse struct {
	PaymentID      uuid.UUID  `json:"payment_id"`
	OrderID        uuid.UUID  `json:"order_id"`
	Method         string     `json:"method"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"`
	TransactionRef *string    `json:"transaction_ref,omitempty"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func FromPaymentModelToResponse(
	payment *model.Payment,
) *PaymentResponse {
	return &PaymentResponse{
		PaymentID:      payment.PaymentID,
		OrderID:        payment.OrderID,
		Method:         payment.Method,
		Amount:         payment.Amount,
		Status:         payment.Status,
		TransactionRef: payment.TransactionRef,
		PaidAt:         payment.PaidAt,
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}
}
