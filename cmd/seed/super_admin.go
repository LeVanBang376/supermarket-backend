package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/config"
	"supermarket-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedSuperAdmin(
	db *gorm.DB,
	cfg *config.Config,
) error {
	var existingUser model.User

	result := db.
		Where("username = ?", cfg.SuperAdminUsername).
		First(&existingUser)

	// Super admin already exists.
	if result.Error == nil {
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf(
			"check super admin: %w",
			result.Error,
		)
	}

	// Hash password before storing it.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.SuperAdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf(
			"hash super admin password: %w",
			err,
		)
	}

	user := &model.User{
		Username:     cfg.SuperAdminUsername,
		PasswordHash: string(passwordHash),
		Email:        cfg.SuperAdminEmail,
		FullName:     cfg.SuperAdminFullName,
		PhoneNumber:  cfg.SuperAdminPhone,
		BranchID:     model.DefaultBranchID,
		RoleID:       model.RoleSuperAdmin,
		PositionID:   model.PositionAdmin,
		Status:       model.UserStatusActive,
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf(
			"create super admin: %w",
			err,
		)
	}

	return nil
}
