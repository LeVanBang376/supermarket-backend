package order_item

import (
	"context"

	"supermarket-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	item *model.OrderItem,
) error {
	return db.
		WithContext(ctx).
		Create(item).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
	skuBarcode string,
) (*model.OrderItem, error) {
	var item model.OrderItem

	err := db.
		WithContext(ctx).
		Preload("SKU").
		Preload("SKU.Unit").
		First(
			&item,
			"order_id = ? AND sku_barcode = ?",
			orderID,
			skuBarcode,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) FindByOrderID(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
) ([]*model.OrderItem, error) {
	var items []*model.OrderItem

	err := db.
		WithContext(ctx).
		Preload("SKU").
		Preload("SKU.Unit").
		Where("order_id = ?", orderID).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	item *model.OrderItem,
) error {
	return db.
		WithContext(ctx).
		Model(&model.OrderItem{}).
		Where(
			"order_id = ? AND sku_barcode = ?",
			item.OrderID,
			item.SKUBarcode,
		).
		Updates(item).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
	skuBarcode string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.OrderItem{},
			"order_id = ? AND sku_barcode = ?",
			orderID,
			skuBarcode,
		).
		Error
}
