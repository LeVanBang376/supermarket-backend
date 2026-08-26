package dto

import (
	"supermarket-backend/internal/model"
	"time"
)

type RoleResponse struct {
	RoleID    string    `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromRoleModelToResponse(role *model.Role) *RoleResponse {
	var r *RoleResponse = &RoleResponse{}

	r.RoleID = role.RoleID
	r.RoleName = role.RoleName

	return r
}
