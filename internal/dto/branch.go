package dto

type CreateBranchRequest struct {
	BranchName string `json:"branch_name" binding:"required"`
	Address    string `json:"address" binding:"required"`
}

type UpdateBranchRequest struct {
	BranchName *string `json:"branch_name"`
	Address    *string `json:"address"`
}

type BranchResponse struct {
	BranchID   string `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Address    string `json:"address"`
}
