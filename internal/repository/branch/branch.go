package branch

import (
	"context"
	"fmt"

	"supermarket-backend/internal/model"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	branch *model.Branch,
) error {
	return db.
		WithContext(ctx).
		Create(branch).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	branchID string,
) (*model.Branch, error) {
	var branch model.Branch

	err := db.
		WithContext(ctx).
		First(&branch, "branch_id = ?", branchID).
		Error

	if err != nil {
		return nil, err
	}

	return &branch, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.Branch, error) {
	var branches []*model.Branch

	// Count total records
	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Branch{}).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.Total = total

	// Calculate total pages
	pagination.TotalPages = int(
		(total + int64(pagination.PerPage) - 1) /
			int64(pagination.PerPage),
	)

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.PerPage

	// Get paginated branches
	if err := db.
		WithContext(ctx).
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&branches).
		Error; err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	branch *model.Branch,
) error {
	return db.
		WithContext(ctx).
		Model(&model.Branch{}).
		Where("branch_id = ?", branch.BranchID).
		Updates(branch).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	branchID string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Branch{},
			"branch_id = ?",
			branchID,
		).
		Error
}

func (r *Repository) GetNextBranchID(
	ctx context.Context,
	db *gorm.DB,
) (string, error) {
	var lastBranchID string

	err := db.
		WithContext(ctx).
		Model(&model.Branch{}).
		Select("branch_id").
		Order("branch_id DESC").
		Limit(1).
		Scan(&lastBranchID).
		Error

	if err != nil {
		return "", err
	}

	if lastBranchID == "" {
		return model.DefaultBranchID, nil
	}

	var number int

	_, err = fmt.Sscanf(lastBranchID, "BR%d", &number)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("BR%04d", number+1), nil
}
