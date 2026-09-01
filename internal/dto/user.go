package dto

import (
	"supermarket-backend/internal/model"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	Email       string `json:"email" binding:"required,email"`
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
		Email:        r.Email,
		FullName:     r.FullName,
		PhoneNumber:  r.PhoneNumber,
		BranchID:     r.BranchID,
		RoleID:       r.RoleID,
		PositionID:   r.PositionID,
	}
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	Email       *string `json:"email" binding:"omitempty,email"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	BranchID    *string `json:"branch_id"`
	RoleID      *string `json:"role_id"`
	PositionID  *string `json:"position_id"`
	Status      *string `json:"status"`
}

type UserResponse struct {
	UserID      uuid.UUID         `json:"user_id"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	FullName    string            `json:"full_name"`
	PhoneNumber string            `json:"phone_number"`
	Branch      *BranchResponse   `json:"branch"`
	Role        *RoleResponse     `json:"role"`
	Position    *PositionResponse `json:"position"`
	Status      string            `json:"status"`
}

func FromUserModelToResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}

	var u *UserResponse = &UserResponse{}

	u.UserID = user.UserID
	u.Username = user.Username
	u.Email = user.Email
	u.FullName = user.FullName
	u.PhoneNumber = user.PhoneNumber
	u.Status = user.Status

	u.Branch = FromBranchModelToResponse(&user.Branch)
	u.Role = FromRoleModelToResponse(&user.Role)
	u.Position = FromPositionModelToResponse(&user.Position)

	return u
}
