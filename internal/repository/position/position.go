package position

import (
	"context"

	"supermarket-backend/internal/model"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByID(
	ctx context.Context,
	positionID string,
) (*model.Position, error) {
	var position model.Position

	err := r.db.
		WithContext(ctx).
		First(&position, "position_id = ?", positionID).
		Error

	if err != nil {
		return nil, err
	}

	return &position, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*model.Position, error) {
	var positions []*model.Position

	// Default pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.PerPage < 1 {
		pagination.PerPage = 10
	}

	// Count total
	var total int64
	if err := r.db.
		WithContext(ctx).
		Model(&model.Position{}).
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

	// Get positions
	if err := r.db.
		WithContext(ctx).
		Order("position_id ASC").
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&positions).
		Error; err != nil {
		return nil, err
	}

	return positions, nil
}
