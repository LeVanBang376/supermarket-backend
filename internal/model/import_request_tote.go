package model

import "time"

type ImportRequestTote struct {
	RequestID   string    `gorm:"column:request_id;type:varchar(7);primaryKey;not null" json:"request_id"`
	ToteBarcode string    `gorm:"column:tote_barcode;type:varchar(30);primaryKey;not null" json:"tote_barcode"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Request ImportRequest `gorm:"foreignKey:RequestID;references:RequestID" json:"request"`
}

func (ImportRequestTote) TableName() string {
	return "import_request_totes"
}
