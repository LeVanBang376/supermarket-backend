package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/config"
	"supermarket-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedAdmin(
	db *gorm.DB,
	cfg *config.Config,
) error {
	var existingUser model.User

	result := db.
		Where("username = ?", cfg.AdminUsername).
		First(&existingUser)

	// Admin already exists.
	if result.Error == nil {
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf(
			"check admin: %w",
			result.Error,
		)
	}

	// Hash password before storing it.
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
		BranchID:     "BR0001",
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

	return nil
}
