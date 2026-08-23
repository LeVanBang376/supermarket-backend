package main

import (
	"supermarket-backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedPositions(db *gorm.DB) error {
	positions := []model.Position{
		{
			PositionID:   model.PositionAdmin,
			PositionName: model.PositionAdminName,
		},
		{
			PositionID:   model.PositionStoreManager,
			PositionName: model.PositionStoreManagerName,
		},
		{
			PositionID:   model.PositionCashier,
			PositionName: model.PositionCashierName,
		},
		{
			PositionID:   model.PositionEmployee,
			PositionName: model.PositionEmployeeName,
		},
		{
			PositionID:   model.PositionRecruiter,
			PositionName: model.PositionRecruiterName,
		},
	}

	return db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&positions).
		Error
}
