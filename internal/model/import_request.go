package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	ImportRequestStatusDraft            = "DRAFT"
	ImportRequestStatusCancelled        = "CANCELLED"
	ImportRequestStatusRequired         = "REQUIRED"
	ImportRequestStatusSupplierReceived = "SUPPLIER_RECEIVED"
	ImportRequestStatusDelivering       = "DELIVERING"
	ImportRequestStatusRejected         = "REJECTED"
	ImportRequestStatusCompleted        = "COMPLETED"
)

type ImportRequest struct {
	RequestID            string     `gorm:"column:request_id;type:varchar(7);primaryKey" json:"request_id"`
	BranchID             string     `gorm:"column:branch_id;type:varchar(6);not null" json:"branch_id"`
	CreatedBy            uuid.UUID  `gorm:"column:created_by;type:uuid;not null" json:"created_by"`
	ExpectedDeliveryAt   time.Time  `gorm:"column:expected_delivery_at;not null" json:"expected_delivery_at"`
	DeliveryLicensePlate string     `gorm:"column:delivery_license_plate;type:varchar(20);not null" json:"delivery_license_plate"`
	Status               string     `gorm:"column:status;type:varchar(20);not null" json:"status"`
	ReceivedBy           *uuid.UUID `gorm:"column:received_by;type:uuid" json:"received_by"`
	CompleteAt           *time.Time `gorm:"column:complete_at" json:"complete_at"`
	CreatedAt            time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Branch   Branch `gorm:"foreignKey:BranchID;references:BranchID" json:"branch"`
	Creator  User   `gorm:"foreignKey:CreatedBy;references:UserID" json:"creator"`
	Receiver *User  `gorm:"foreignKey:ReceivedBy;references:UserID" json:"receiver"`

	Totes    []ImportRequestTote    `gorm:"foreignKey:RequestID;references:RequestID" json:"totes"`
	Products []ImportRequestProduct `gorm:"foreignKey:RequestID;references:RequestID" json:"products"`
}

func (ImportRequest) TableName() string {
	return "import_requests"
}
