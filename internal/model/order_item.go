package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	OrderID    uuid.UUID `gorm:"column:order_id;type:uuid;primaryKey" json:"order_id"`
	SKUBarcode string    `gorm:"column:sku_barcode;type:varchar(30);primaryKey" json:"sku_barcode"`

	Quantity       float64 `gorm:"column:quantity;type:numeric(10,3);not null" json:"quantity"`
	UnitPrice      float64 `gorm:"column:unit_price;type:numeric(12,2);not null" json:"unit_price"`
	Subtotal       float64 `gorm:"column:subtotal;type:numeric(12,2);not null" json:"subtotal"`
	DiscountAmount float64 `gorm:"column:discount_amount;type:numeric(12,2);not null;default:0" json:"discount_amount"`

	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Order Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order"`
	SKU   SKU   `gorm:"foreignKey:SKUBarcode;references:SKUBarcode" json:"sku"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
