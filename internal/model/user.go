package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	Username     string    `gorm:"column:username;type:varchar(50);unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null" json:"-"`
	FullName     string    `gorm:"column:full_name;type:varchar(30);not null"`
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(15);not null"`
	BranchID     string    `gorm:"column:branch_id;type:varchar(6);not null"`
	RoleID       string    `gorm:"column:role_id;type:varchar(10);not null"`
	PositionID   string    `gorm:"column:position_id;type:varchar(10);not null"`
	Status       string    `gorm:"column:status;type:varchar(20);not null;default:ACTIVE"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()"`

	// Associations
	Branch   Branch   `gorm:"foreignKey:BranchID;references:BranchID"`
	Role     Role     `gorm:"foreignKey:RoleID;references:RoleID"`
	Position Position `gorm:"foreignKey:PositionID;references:PositionID"`
}

func (User) TableName() string {
	return "users"
}
