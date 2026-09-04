package sku

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

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	sku *model.SKU,
) error {
	return db.
		WithContext(ctx).
		Create(sku).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	skuBarcode string,
) (*model.SKU, error) {
	var sku model.SKU

	err := db.
		WithContext(ctx).
		Preload("Brand").
		Preload("Unit").
		First(&sku, "sku_barcode = ?", skuBarcode).
		Error

	if err != nil {
		return nil, err
	}

	return &sku, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.SKU, error) {
	var skus []*model.SKU

	var total int64

	if err := db.
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

	if err := db.
		WithContext(ctx).
		Preload("Brand").
		Preload("Unit").
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
	db *gorm.DB,
	sku *model.SKU,
) error {
	return db.
		WithContext(ctx).
		Model(&model.SKU{}).
		Where("sku_barcode = ?", sku.SKUBarcode).
		Updates(sku).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	skuBarcode string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.SKU{},
			"sku_barcode = ?",
			skuBarcode,
		).
		Error
}
