package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/config"
	"supermarket-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUsers(
	db *gorm.DB,
	cfg *config.Config,
) error {
	// ============================================
	// Admin
	// ============================================

	var existingUser model.User

	result := db.
		Where("username = ?", cfg.AdminUsername).
		First(&existingUser)

	// Admin already exists.
	if result.Error == nil {
		// Do nothing.
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check admin: %w",
				result.Error,
			)
		}

		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(cfg.AdminPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return fmt.Errorf(
				"hash admin password: %w",
				err,
			)
		}

		user := &model.User{
			Username:     cfg.AdminUsername,
			PasswordHash: string(passwordHash),
			Email:        cfg.AdminEmail,
			FullName:     cfg.AdminFullName,
			PhoneNumber:  cfg.AdminPhone,
			BranchID:     model.DefaultBranchID,
			RoleID:       model.RoleAdmin,
			PositionID:   model.PositionAdmin,
			Status:       model.UserStatusActive,
		}

		if err := db.Create(user).Error; err != nil {
			return fmt.Errorf(
				"create admin: %w",
				err,
			)
		}
	}

	// ============================================
	// Demo users
	// ============================================

	demoPassword := "123456"

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(demoPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf(
			"hash demo password: %w",
			err,
		)
	}

	demoUsers := []model.User{
		{
			Username:     "sm_staff",
			PasswordHash: string(passwordHash),
			Email:        "employee01@example.com",
			FullName:     "Nguyen Van An",
			PhoneNumber:  "0900000001",
			BranchID:     "BR0001",
			RoleID:       model.RoleEmployee,
			PositionID:   model.PositionEmployee,
			Status:       model.UserStatusActive,
		},
		{
			Username:     "sm_manager",
			PasswordHash: string(passwordHash),
			Email:        "manager01@example.com",
			FullName:     "Tran Van Binh",
			PhoneNumber:  "0900000002",
			BranchID:     "BR0002",
			RoleID:       model.RoleManager,
			PositionID:   model.PositionStoreManager,
			Status:       model.UserStatusActive,
		},
		{
			Username:     "sm_cashier",
			PasswordHash: string(passwordHash),
			Email:        "cashier01@example.com",
			FullName:     "Le Thi Chi",
			PhoneNumber:  "0900000003",
			BranchID:     "BR0002",
			RoleID:       model.RoleEmployee,
			PositionID:   model.PositionCashier,
			Status:       model.UserStatusActive,
		},
		{
			Username:     "hr",
			PasswordHash: string(passwordHash),
			Email:        "hr01@example.com",
			FullName:     "Pham Thi Dung",
			PhoneNumber:  "0900000004",
			BranchID:     "BR0001",
			RoleID:       model.RoleHR,
			PositionID:   model.PositionRecruiter,
			Status:       model.UserStatusActive,
		},
		{
			Username:     "supplier_staff",
			PasswordHash: string(passwordHash),
			Email:        "supplier01@example.com",
			FullName:     "Hoang Van Nam",
			PhoneNumber:  "0900000005",
			BranchID:     "BR0003",
			RoleID:       model.RoleSupplier,
			PositionID:   model.PositionEmployee,
			Status:       model.UserStatusActive,
		},
	}

	for _, user := range demoUsers {
		var existingUser model.User

		result := db.
			Where("username = ?", user.Username).
			First(&existingUser)

		if result.Error == nil {
			// User already exists.
			continue
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check demo user %s: %w",
				user.Username,
				result.Error,
			)
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf(
				"create demo user %s: %w",
				user.Username,
				err,
			)
		}
	}

	return nil
}
