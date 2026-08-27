package main

import (
	"fmt"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedStocks(db *gorm.DB) error {
	stocks := []model.Stock{
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970011",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970028",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970035",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970042",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970059",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970066",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970073",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970080",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970097",
			Quantity:   0,
		},
		{
			BranchID:   "BR0002",
			SKUBarcode: "8938505970103",
			Quantity:   0,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&stocks).Error; err != nil {
		return fmt.Errorf("create stocks: %w", err)
	}

	return nil
}
