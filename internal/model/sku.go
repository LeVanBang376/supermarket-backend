package model

import "time"

type SKU struct {
	SKUBarcode    string    `gorm:"column:sku_barcode;type:varchar(30);primaryKey"`
	SKUName       string    `gorm:"column:sku_name;type:varchar(50);not null"`
	BrandID       string    `gorm:"column:brand_id;type:varchar(5);not null"`
	UnitID        string    `gorm:"column:unit_id;type:varchar(5);not null"`
	ShelfLifeDays int       `gorm:"column:shelf_life_days;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()"`

	Brand Brand `gorm:"foreignKey:BrandID;references:BrandID"`
	Unit  Unit  `gorm:"foreignKey:UnitID;references:UnitID"`
}

func (SKU) TableName() string {
	return "skus"
}
