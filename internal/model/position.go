package model

import "time"

type Position struct {
	PositionID   string    `gorm:"column:position_id;type:varchar(10);primaryKey"`
	PositionName string    `gorm:"column:position_name;type:varchar(30);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Position) TableName() string {
	return "positions"
}
