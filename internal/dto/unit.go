package dto

import "supermarket-backend/internal/model"

type CreateUnitRequest struct {
	UnitName string `json:"unit_name" binding:"required,max=30"`
}

type UpdateUnitRequest struct {
	UnitName *string `json:"unit_name" binding:"omitempty,max=30"`
}

type UnitResponse struct {
	UnitID   string `json:"unit_id"`
	UnitName string `json:"unit_name"`
}

func FromUnitModelToResponse(unit *model.Unit) *UnitResponse {
	var u *UnitResponse = &UnitResponse{}

	u.UnitID = unit.UnitID
	u.UnitName = unit.UnitName

	return u
}
