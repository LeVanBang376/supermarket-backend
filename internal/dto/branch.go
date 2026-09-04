package dto

import "supermarket-backend/internal/model"

type CreateBranchRequest struct {
	BranchName string           `json:"branch_name" binding:"required,max=100"`
	Address    string           `json:"address" binding:"required,max=150"`
	Type       model.BranchType `json:"type" binding:"required,oneof=HEADQUARTERS SUPERMARKET SUPPLIER"`
}

type UpdateBranchRequest struct {
	BranchName *string           `json:"branch_name" binding:"omitempty,max=100"`
	Address    *string           `json:"address" binding:"omitempty,max=150"`
	Type       *model.BranchType `json:"type" binding:"omitempty,oneof=HEADQUARTERS SUPERMARKET SUPPLIER"`
}

type BranchResponse struct {
	BranchID   string           `json:"branch_id"`
	BranchName string           `json:"branch_name"`
	Address    string           `json:"address"`
	Type       model.BranchType `json:"type"`
}

func FromBranchModelToResponse(branch *model.Branch) *BranchResponse {
	return &BranchResponse{
		BranchID:   branch.BranchID,
		BranchName: branch.BranchName,
		Address:    branch.Address,
		Type:       branch.Type,
	}
}
