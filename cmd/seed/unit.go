package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

func seedUnits(db *gorm.DB) error {
	units := []model.Unit{
		{
			UnitID:   "U0001",
			UnitName: "Bó",
		},
		{
			UnitID:   "U0002",
			UnitName: "Khay",
		},
		{
			UnitID:   "U0003",
			UnitName: "Túi",
		},
		{
			UnitID:   "U0004",
			UnitName: "Bắp",
		},
		{
			UnitID:   "U0005",
			UnitName: "Quả",
		},
		{
			UnitID:   "U0006",
			UnitName: "Chai",
		},
		{
			UnitID:   "U0007",
			UnitName: "Cái",
		},
		{
			UnitID:   "U0008",
			UnitName: "Hộp",
		},
		{
			UnitID:   "U0009",
			UnitName: "Lốc",
		},
		{
			UnitID:   "U0010",
			UnitName: "Gói",
		},
	}

	for _, unit := range units {
		var existingUnit model.Unit

		result := db.
			Where("unit_id = ?", unit.UnitID).
			First(&existingUnit)

		if result.Error == nil {
			continue
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check unit %s: %w",
				unit.UnitID,
				result.Error,
			)
		}

		if err := db.Create(&unit).Error; err != nil {
			return fmt.Errorf(
				"create unit %s: %w",
				unit.UnitID,
				err,
			)
		}
	}

	return nil
}
