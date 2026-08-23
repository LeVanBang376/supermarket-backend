package model

import "time"

const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleManager    = "MANAGER"
	RoleEmployee   = "EMPLOYEE"
	RoleHR         = "HR"
)

const (
	RoleSuperAdminName = "Super Admin"
	RoleAdminName      = "Admin"
	RoleManagerName    = "Manager"
	RoleEmployeeName   = "Employee"
	RoleHRName         = "HR"
)

type Role struct {
	RoleID    string    `gorm:"column:role_id;type:varchar(20);primaryKey" json:"role_id"`
	RoleName  string    `gorm:"column:role_name;type:varchar(30);not null" json:"role_name"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
