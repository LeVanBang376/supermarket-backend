package dto

import "supermarket-backend/internal/model"

type CreateBranchRequest struct {
	BranchName string `json:"branch_name" binding:"required,max=100"`
	Address    string `json:"address" binding:"required,max=150"`
}

type UpdateBranchRequest struct {
	BranchName *string `json:"branch_name" binding:"omitempty,max=100"`
	Address    *string `json:"address" binding:"omitempty,max=150"`
}

type BranchResponse struct {
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Address    string `json:"address"`
}

func FromBranchModelToResponse(branch *model.Branch) *BranchResponse {
	var b *BranchResponse = &BranchResponse{}

	b.BranchID = branch.BranchID
	b.BranchName = branch.BranchName
	b.Address = branch.Address

	return b
}
