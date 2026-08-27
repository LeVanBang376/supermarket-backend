package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

func seedBrands(db *gorm.DB) error {
	brands := []model.Brand{
		{
			BrandID:   "B0001",
			BrandName: "VinEco",
		},
		{
			BrandID:   "B0002",
			BrandName: "Dalat GAP",
		},
		{
			BrandID:   "B0003",
			BrandName: "Organica",
		},
		{
			BrandID:   "B0004",
			BrandName: "BigGreen",
		},
	}

	for _, brand := range brands {
		var existingBrand model.Brand

		result := db.
			Where("brand_id = ?", brand.BrandID).
			First(&existingBrand)

		if result.Error == nil {
			continue
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check brand %s: %w",
				brand.BrandID,
				result.Error,
			)
		}

		if err := db.Create(&brand).Error; err != nil {
			return fmt.Errorf(
				"create brand %s: %w",
				brand.BrandID,
				err,
			)
		}
	}

	return nil
}
