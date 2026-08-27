package main

import (
	"errors"
	"fmt"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

func seedSKUs(db *gorm.DB) error {
	skus := []model.SKU{
		{
			SKUBarcode:    "8938505970011",
			SKUName:       "Rau muống VietGAP",
			BrandID:       "B0001",
			UnitID:        "U0001",
			ShelfLifeDays: 3,
		},
		{
			SKUBarcode:    "8938505970028",
			SKUName:       "Cải ngọt VietGAP",
			BrandID:       "B0002",
			UnitID:        "U0001",
			ShelfLifeDays: 4,
		},
		{
			SKUBarcode:    "8938505970035",
			SKUName:       "Cải xanh VietGAP",
			BrandID:       "B0002",
			UnitID:        "U0001",
			ShelfLifeDays: 4,
		},
		{
			SKUBarcode:    "8938505970042",
			SKUName:       "Mồng tơi VietGAP",
			BrandID:       "B0003",
			UnitID:        "U0001",
			ShelfLifeDays: 3,
		},
		{
			SKUBarcode:    "8938505970059",
			SKUName:       "Cà chua VietGAP",
			BrandID:       "B0003",
			UnitID:        "U0005",
			ShelfLifeDays: 7,
		},
		{
			SKUBarcode:    "8938505970066",
			SKUName:       "Dưa leo VietGAP",
			BrandID:       "B0004",
			UnitID:        "U0005",
			ShelfLifeDays: 7,
		},
		{
			SKUBarcode:    "8938505970073",
			SKUName:       "Ớt đỏ VietGAP",
			BrandID:       "B0004",
			UnitID:        "U0005",
			ShelfLifeDays: 10,
		},
		{
			SKUBarcode:    "8938505970080",
			SKUName:       "Bắp cải VietGAP",
			BrandID:       "B0001",
			UnitID:        "U0004",
			ShelfLifeDays: 14,
		},
		{
			SKUBarcode:    "8938505970097",
			SKUName:       "Bí đỏ VietGAP",
			BrandID:       "B0001",
			UnitID:        "U0005",
			ShelfLifeDays: 21,
		},
		{
			SKUBarcode:    "8938505970103",
			SKUName:       "Su su VietGAP",
			BrandID:       "B0002",
			UnitID:        "U0005",
			ShelfLifeDays: 10,
		},
	}

	for _, sku := range skus {
		var existingSKU model.SKU

		result := db.
			Where("sku_barcode = ?", sku.SKUBarcode).
			First(&existingSKU)

		if result.Error == nil {
			continue
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf(
				"check sku %s: %w",
				sku.SKUBarcode,
				result.Error,
			)
		}

		if err := db.Create(&sku).Error; err != nil {
			return fmt.Errorf(
				"create sku %s: %w",
				sku.SKUBarcode,
				err,
			)
		}
	}

	return nil
}
