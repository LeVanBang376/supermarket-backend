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
			BranchID:   model.DefaultBranchID,
			BranchName: "Chi nhánh trung tâm",
			Address:    "Hồ Chí Minh",
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
