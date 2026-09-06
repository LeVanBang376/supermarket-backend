package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	PaymentMethodCash = "CASH"
	PaymentMethodQR   = "QR"
	PaymentMethodCard = "CARD"

	PaymentStatusPending   = "PENDING"
	PaymentStatusSuccess   = "SUCCESS"
	PaymentStatusFailed    = "FAILED"
	PaymentStatusCancelled = "CANCELLED"
)

type Payment struct {
	PaymentID uuid.UUID `gorm:"column:payment_id;type:uuid;primaryKey" json:"payment_id"`
	OrderID   uuid.UUID `gorm:"column:order_id;type:uuid;not null" json:"order_id"`

	Method string  `gorm:"column:method;type:varchar(20);not null" json:"method"`
	Amount float64 `gorm:"column:amount;type:numeric(12,2);not null" json:"amount"`
	Status string  `gorm:"column:status;type:varchar(20);not null" json:"status"`

	TransactionRef *string    `gorm:"column:transaction_ref;type:varchar(100)" json:"transaction_ref"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at"`

	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Order Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order"`
}

func (Payment) TableName() string {
	return "payments"
}
