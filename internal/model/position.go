package model

import "time"

const (
	PositionAdmin        = "ADMIN"
	PositionStoreManager = "STORE_MANAGER"
	PositionCashier      = "CASHIER"
	PositionEmployee     = "EMPLOYEE"
	PositionRecruiter    = "RECRUITER"
)

const (
	PositionAdminName        = "Quản lý vận hành"
	PositionStoreManagerName = "Cửa hàng trưởng"
	PositionCashierName      = "Thu ngân"
	PositionEmployeeName     = "Nhân viên"
	PositionRecruiterName    = "Nhân viên tuyển dụng"
)

type Position struct {
	PositionID   string    `gorm:"column:position_id;type:varchar(20);primaryKey"`
	PositionName string    `gorm:"column:position_name;type:varchar(30);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()"`
}

func (Position) TableName() string {
	return "positions"
}
