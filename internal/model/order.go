package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusOpen      = "OPEN"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCancelled = "CANCELLED"
)

type Order struct {
	OrderID   uuid.UUID `gorm:"column:order_id;type:uuid;primaryKey" json:"order_id"`
	BranchID  string    `gorm:"column:branch_id;type:varchar(6);not null" json:"branch_id"`
	CashierID uuid.UUID `gorm:"column:cashier_id;type:uuid;not null" json:"cashier_id"`
	Status    string    `gorm:"column:status;type:varchar(20);not null" json:"status"`

	Subtotal       float64 `gorm:"column:subtotal;type:numeric(12,2);not null;default:0" json:"subtotal"`
	DiscountAmount float64 `gorm:"column:discount_amount;type:numeric(12,2);not null;default:0" json:"discount_amount"`
	TotalAmount    float64 `gorm:"column:total_amount;type:numeric(12,2);not null;default:0" json:"total_amount"`

	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Branch  Branch `gorm:"foreignKey:BranchID;references:BranchID" json:"branch"`
	Cashier User   `gorm:"foreignKey:CashierID;references:UserID" json:"cashier"`

	Items    []OrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"items"`
	Payments []Payment   `gorm:"foreignKey:OrderID;references:OrderID" json:"payments"`
}

func (Order) TableName() string {
	return "orders"
}
