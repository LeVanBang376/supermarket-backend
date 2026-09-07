package dto

import (
	"time"

	"github.com/google/uuid"

	"supermarket-backend/internal/model"
)

type FindAllOrdersQuery struct {
	BranchID *string `form:"branch_id"`
	Status   *string `form:"status"`
}

type CreateOrderRequest struct {
	BranchID string `json:"branch_id" binding:"required,max=6"`
}

type UpdateOrderRequest struct {
	Status *string `json:"status" binding:"omitempty,max=20"`
}

type OrderResponse struct {
	OrderID        uuid.UUID `json:"order_id"`
	BranchID       string    `json:"branch_id"`
	CashierID      uuid.UUID `json:"cashier_id"`
	Status         string    `json:"status"`
	Subtotal       float64   `json:"subtotal"`
	DiscountAmount float64   `json:"discount_amount"`
	TotalAmount    float64   `json:"total_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Items    []*OrderItemResponse `json:"items"`
	Payments []*PaymentResponse `json:"payments"`
}

func FromOrderModelToResponse(order *model.Order) *OrderResponse {
	// items := make(
	// 	[]*OrderItemResponse,
	// 	0,
	// 	len(order.Items),
	// )

	// for _, item := range order.Items {
	// 	items = append(
	// 		items,
	// 		FromOrderItemModelToResponse(&item),
	// 	)
	// }

	payments := make(
		[]*PaymentResponse,
		0,
		len(order.Payments),
	)

	for _, payment := range order.Payments {
		payments = append(
			payments,
			FromPaymentModelToResponse(&payment),
		)
	}

	return &OrderResponse{
		OrderID:        order.OrderID,
		BranchID:       order.BranchID,
		CashierID:      order.CashierID,
		Status:         order.Status,
		Subtotal:       order.Subtotal,
		DiscountAmount: order.DiscountAmount,
		TotalAmount:    order.TotalAmount,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
		// Items:          items,
		Payments: payments,
	}
}
