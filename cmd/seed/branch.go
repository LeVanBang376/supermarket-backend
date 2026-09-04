package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

func seedBranches(db *gorm.DB) error {
	branches := []model.Branch{
		{
			BranchID:   "BR0001",
			BranchName: "Chi nhánh trung tâm",
			Address:    "Hồ Chí Minh",
			Type:       model.BranchTypeHeadquarters,
		},
		{
			BranchID:   "BR0002",
			BranchName: "Siêu thị Quận 1",
			Address:    "Quận 1, Hồ Chí Minh",
			Type:       model.BranchTypeSupermarket,
		},
		{
			BranchID:   "BR0002",
			BranchName: "Nhà cung cấp TP.Hồ Chí Minh",
			Address:    "Quận 7, Hồ Chí Minh",
			Type:       model.BranchTypeSupplier,
		},
	}

	for _, branch := range branches {
		var existingBranch model.Branch

		result := db.
			Where("branch_id = ?", branch.BranchID).
			First(&existingBranch)

		if result.Error == nil {
			continue
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check branch %s: %w",
				branch.BranchID,
				result.Error,
			)
		}

		if err := db.Create(&branch).Error; err != nil {
			return fmt.Errorf(
				"create branch %s: %w",
				branch.BranchID,
				err,
			)
		}
	}

	return nil
}
