package model

import "time"

type Brand struct {
	BrandID   string    `gorm:"column:brand_id;type:varchar(5);primaryKey"`
	BrandName string    `gorm:"column:brand_name;type:varchar(30);not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Brand) TableName() string {
	return "brands"
}
