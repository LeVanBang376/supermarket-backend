package dto

import (
	"time"

	"github.com/google/uuid"

	"supermarket-backend/internal/model"
)

type AddOrderItemRequest struct {
	SKUBarcode string  `json:"sku_barcode" binding:"required,max=30"`
	Quantity   float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderItemRequest struct {
	Quantity *float64 `json:"quantity" binding:"omitempty,gt=0"`
}

type OrderItemResponse struct {
	OrderID        uuid.UUID `json:"order_id"`
	SKUBarcode     string    `json:"sku_barcode"`
	SKUName        string    `json:"sku_name"`
	UnitID         string    `json:"unit_id"`
	UnitName       string    `json:"unit_name"`
	Quantity       float64   `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	Subtotal       float64   `json:"subtotal"`
	DiscountAmount float64   `json:"discount_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func FromOrderItemModelToResponse(
	item *model.OrderItem,
) *OrderItemResponse {
	return &OrderItemResponse{
		OrderID:        item.OrderID,
		SKUBarcode:     item.SKUBarcode,
		SKUName:        item.SKU.SKUName,
		UnitID:         item.SKU.UnitID,
		UnitName:       item.SKU.Unit.UnitName,
		Quantity:       item.Quantity,
		UnitPrice:      item.UnitPrice,
		Subtotal:       item.Subtotal,
		DiscountAmount: item.DiscountAmount,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}
