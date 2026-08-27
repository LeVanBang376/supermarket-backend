package sku

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

func (r *Repository) Create(
	ctx context.Context,
	sku *model.SKU,
) error {
	return r.db.
		WithContext(ctx).
		Create(sku).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	skuBarcode string,
) (*model.SKU, error) {
	var sku model.SKU

	err := r.db.
		WithContext(ctx).
		First(&sku, "sku_barcode = ?", skuBarcode).
		Error

	if err != nil {
		return nil, err
	}

	return &sku, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*model.SKU, error) {
	var skus []*model.SKU

	var total int64
	if err := r.db.
		WithContext(ctx).
		Model(&model.SKU{}).
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

	if err := r.db.
		WithContext(ctx).
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&skus).
		Error; err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *Repository) Update(
	ctx context.Context,
	sku *model.SKU,
) error {
	return r.db.
		WithContext(ctx).
		Model(&model.SKU{}).
		Where("sku_barcode = ?", sku.SKUBarcode).
		Updates(sku).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	skuBarcode string,
) error {
	return r.db.
		WithContext(ctx).
		Delete(
			&model.SKU{},
			"sku_barcode = ?",
			skuBarcode,
		).
		Error
}
