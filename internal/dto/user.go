package dto

import (
	"supermarket-backend/internal/model"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	BranchID    string `json:"branch_id" binding:"required"`
	RoleID      string `json:"role_id" binding:"required"`
	PositionID  string `json:"position_id" binding:"required"`
}

func (r CreateUserRequest) ToUserModel(passwordHash string) *model.User {
	return &model.User{
		Username:     r.Username,
		PasswordHash: passwordHash,
		FullName:     r.FullName,
		PhoneNumber:  r.PhoneNumber,
		BranchID:     r.BranchID,
		RoleID:       r.RoleID,
		PositionID:   r.PositionID,
	}
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	BranchID    *string `json:"branch_id"`
	RoleID      *string `json:"role_id"`
	PositionID  *string `json:"position_id"`
	Status      *string `json:"status"`
}

type UserResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	BranchID    string    `json:"branch_id"`
	RoleID      string    `json:"role_id"`
	PositionID  string    `json:"position_id"`
	Status      string    `json:"status"`
}
