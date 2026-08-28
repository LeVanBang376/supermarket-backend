package role

import (
	"context"

	"supermarket-backend/internal/model"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	roleID string,
) (*model.Role, error) {
	var role model.Role

	err := db.
		WithContext(ctx).
		First(&role, "role_id = ?", roleID).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.Role, error) {
	var roles []*model.Role

	// Count total records
	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Role{}).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.Total = total

	// Calculate total pages
	if pagination.PerPage > 0 {
		pagination.TotalPages = int(
			(total + int64(pagination.PerPage) - 1) /
				int64(pagination.PerPage),
		)
	}

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.PerPage

	// Get paginated records
	if err := db.
		WithContext(ctx).
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&roles).
		Error; err != nil {
		return nil, err
	}

	return roles, nil
}
