package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusActive   = "ACTIVE"
	UserStatusInactive = "INACTIVE"
)

type User struct {
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey;default:gen_random_uuid()" json:"user_id"`
	Username     string    `gorm:"column:username;type:varchar(50);unique;not null" json:"username"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null" json:"-"`
	Email        string    `gorm:"column:email;type:varchar(254);unique;not null" json:"email"`
	FullName     string    `gorm:"column:full_name;type:varchar(30);not null" json:"full_name"`
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(15);not null" json:"phone_number"`
	BranchID     string    `gorm:"column:branch_id;type:varchar(6);not null" json:"branch_id"`
	RoleID       string    `gorm:"column:role_id;type:varchar(10);not null" json:"role_id"`
	PositionID   string    `gorm:"column:position_id;type:varchar(10);not null" json:"position_id"`
	Status       string    `gorm:"column:status;type:varchar(20);not null;default:ACTIVE" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`

	// Associations
	Branch   Branch   `gorm:"foreignKey:BranchID;references:BranchID" json:"branch"`
	Role     Role     `gorm:"foreignKey:RoleID;references:RoleID" json:"role"`
	Position Position `gorm:"foreignKey:PositionID;references:PositionID" json:"position"`
}

func (User) TableName() string {
	return "users"
}
