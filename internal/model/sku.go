package model

import "time"

type SKU struct {
	SKUBarcode    string    `gorm:"column:sku_barcode;type:varchar(30);primaryKey" json:"sku_barcode"`
	SKUName       string    `gorm:"column:sku_name;type:varchar(50);not null" json:"sku_name"`
	BrandID       string    `gorm:"column:brand_id;type:varchar(5);not null" json:"brand_id"`
	UnitID        string    `gorm:"column:unit_id;type:varchar(5);not null" json:"unit_id"`
	UnitPrice     float64   `gorm:"column:unit_price;type:numeric(12,2);not null" json:"unit_price"`
	ShelfLifeDays int       `gorm:"column:shelf_life_days;not null" json:"shelf_life_days"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Brand Brand `gorm:"foreignKey:BrandID;references:BrandID" json:"brand"`
	Unit  Unit  `gorm:"foreignKey:UnitID;references:UnitID" json:"unit"`
}

func (SKU) TableName() string {
	return "skus"
}
