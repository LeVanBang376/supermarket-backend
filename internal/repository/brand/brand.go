package brand

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
	brand *model.Brand,
) error {
	return db.
		WithContext(ctx).
		Create(brand).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	brandID string,
) (*model.Brand, error) {
	var brand model.Brand

	err := db.
		WithContext(ctx).
		First(&brand, "brand_id = ?", brandID).
		Error

	if err != nil {
		return nil, err
	}

	return &brand, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.Brand, error) {
	var brands []*model.Brand

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Brand{}).
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
		Find(&brands).
		Error; err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	brand *model.Brand,
) error {
	return db.
		WithContext(ctx).
		Model(&model.Brand{}).
		Where("brand_id = ?", brand.BrandID).
		Updates(brand).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	brandID string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Brand{},
			"brand_id = ?",
			brandID,
		).
		Error
}

func (r *Repository) GetNextBrandID(
	ctx context.Context,
	db *gorm.DB,
) (string, error) {
	var lastBrandID string

	err := db.
		WithContext(ctx).
		Model(&model.Brand{}).
		Select("brand_id").
		Order("brand_id DESC").
		Limit(1).
		Scan(&lastBrandID).
		Error

	if err != nil {
		return "", err
	}

	if lastBrandID == "" {
		return "BR001", nil
	}

	var number int

	_, err = fmt.Sscanf(lastBrandID, "BR%d", &number)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("BR%03d", number+1), nil
}
