package model

import "time"

const DefaultBranchID = "BR0001"

type Branch struct {
	BranchID   string    `gorm:"column:branch_id;type:varchar(6);primaryKey"`
	BranchName string    `gorm:"column:branch_name;type:varchar(100);not null"`
	Address    string    `gorm:"column:address;type:varchar(150);not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:now()"`

	Users []User `gorm:"foreignKey:BranchID;references:BranchID"`
}

func (Branch) TableName() string {
	return "branches"
}
