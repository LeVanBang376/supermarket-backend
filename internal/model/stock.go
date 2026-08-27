package model

import "time"

type Stock struct {
	BranchID   string    `gorm:"column:branch_id;type:varchar(6);primaryKey"`
	SKUBarcode string    `gorm:"column:sku_barcode;type:varchar(30);primaryKey"`
	Quantity   int       `gorm:"column:quantity;not null;default:0"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()"`

	Branch Branch `gorm:"foreignKey:BranchID;references:BranchID"`
	SKU    SKU    `gorm:"foreignKey:SKUBarcode;references:SKUBarcode"`
}

func (Stock) TableName() string {
	return "stocks"
}
