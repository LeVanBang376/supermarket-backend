package model

import "time"

type Role struct {
	RoleID    string    `gorm:"column:role_id;primaryKey;size:10" json:"role_id"`
	RoleName  string    `gorm:"column:role_name;size:30;not null" json:"role_name"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
