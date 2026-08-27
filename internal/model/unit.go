package model

import "time"

type Unit struct {
	UnitID    string    `gorm:"column:unit_id;type:varchar(5);primaryKey"`
	UnitName  string    `gorm:"column:unit_name;type:varchar(30);not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Unit) TableName() string {
	return "units"
}
