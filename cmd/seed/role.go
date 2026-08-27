package main

import (
	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

func seedRoles(db *gorm.DB) error {
	roles := []model.Role{
		{
			RoleID:   model.RoleAdmin,
			RoleName: model.RoleAdminName,
		},
		{
			RoleID:   model.RoleManager,
			RoleName: model.RoleManagerName,
		},
		{
			RoleID:   model.RoleEmployee,
			RoleName: model.RoleEmployeeName,
		},
		{
			RoleID:   model.RoleHR,
			RoleName: model.RoleHRName,
		},
	}

	for _, role := range roles {
		var existing model.Role

		err := db.
			Where("role_id = ?", role.RoleID).
			First(&existing).
			Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					return err
				}

				continue
			}

			return err
		}
	}

	return nil
}
