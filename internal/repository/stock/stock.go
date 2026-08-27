package stock

import (
	"context"
	"errors"

	"supermarket-backend/internal/model"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByBranchAndSKU(
	ctx context.Context,
	branchID string,
	skuBarcode string,
) (*model.Stock, error) {
	var stock model.Stock

	err := r.db.
		WithContext(ctx).
		First(
			&stock,
			"branch_id = ? AND sku_barcode = ?",
			branchID,
			skuBarcode,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*model.Stock, error) {
	var stocks []*model.Stock

	var total int64
	if err := r.db.
		WithContext(ctx).
		Model(&model.Stock{}).
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
		Find(&stocks).
		Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *Repository) IncreaseQuantity(
	ctx context.Context,
	branchID string,
	skuBarcode string,
	quantity int,
) error {
	return r.db.
		WithContext(ctx).
		Model(&model.Stock{}).
		Where(
			"branch_id = ? AND sku_barcode = ?",
			branchID,
			skuBarcode,
		).
		UpdateColumn(
			"quantity",
			gorm.Expr("quantity + ?", quantity),
		).
		Error
}

func (r *Repository) DecreaseQuantity(
	ctx context.Context,
	branchID string,
	skuBarcode string,
	quantity int,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&model.Stock{}).
		Where(
			"branch_id = ? AND sku_barcode = ? AND quantity >= ?",
			branchID,
			skuBarcode,
			quantity,
		).
		UpdateColumn(
			"quantity",
			gorm.Expr("quantity - ?", quantity),
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}

	return nil
}

func (r *Repository) SetQuantity(
	ctx context.Context,
	branchID string,
	skuBarcode string,
	quantity int,
) error {
	return r.db.
		WithContext(ctx).
		Model(&model.Stock{}).
		Where(
			"branch_id = ? AND sku_barcode = ?",
			branchID,
			skuBarcode,
		).
		UpdateColumn(
			"quantity",
			quantity,
		).
		Error
}
