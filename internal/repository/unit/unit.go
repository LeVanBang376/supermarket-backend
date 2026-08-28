package unit

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
	unit *model.Unit,
) error {
	return db.
		WithContext(ctx).
		Create(unit).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	unitID string,
) (*model.Unit, error) {
	var unit model.Unit

	err := db.
		WithContext(ctx).
		First(&unit, "unit_id = ?", unitID).
		Error

	if err != nil {
		return nil, err
	}

	return &unit, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.Unit, error) {
	var units []*model.Unit

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Unit{}).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.Total = total

	pagination.TotalPages = int(
		(total + int64(pagination.PerPage) - 1) /
			int64(pagination.PerPage),
	)

	offset := (pagination.Page - 1) * pagination.PerPage

	if err := db.
		WithContext(ctx).
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&units).
		Error; err != nil {
		return nil, err
	}

	return units, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	unit *model.Unit,
) error {
	return db.
		WithContext(ctx).
		Model(&model.Unit{}).
		Where("unit_id = ?", unit.UnitID).
		Updates(unit).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	unitID string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Unit{},
			"unit_id = ?",
			unitID,
		).
		Error
}

func (r *Repository) GetNextUnitID(
	ctx context.Context,
	db *gorm.DB,
) (string, error) {
	var lastUnitID string

	err := db.
		WithContext(ctx).
		Model(&model.Unit{}).
		Select("unit_id").
		Order("unit_id DESC").
		Limit(1).
		Scan(&lastUnitID).
		Error

	if err != nil {
		return "", err
	}

	if lastUnitID == "" {
		return "UN001", nil
	}

	var number int

	_, err = fmt.Sscanf(lastUnitID, "UN%d", &number)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("UN%03d", number+1), nil
}
