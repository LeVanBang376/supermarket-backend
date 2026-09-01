package import_request_product

import (
	"context"

	"supermarket-backend/internal/model"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	product *model.ImportRequestProduct,
) error {
	return db.
		WithContext(ctx).
		Create(product).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
	skuBarcode string,
) (*model.ImportRequestProduct, error) {
	var product model.ImportRequestProduct

	err := db.
		WithContext(ctx).
		Preload("SKU").
		Preload("Tote").
		First(
			&product,
			"request_id = ? AND sku_barcode = ?",
			requestID,
			skuBarcode,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repository) FindByRequestID(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
) ([]*model.ImportRequestProduct, error) {
	var products []*model.ImportRequestProduct

	err := db.
		WithContext(ctx).
		Preload("SKU").
		Preload("Tote").
		Where("request_id = ?", requestID).
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) FindByRequestIDAndToteBarcode(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
	toteBarcode string,
) ([]*model.ImportRequestProduct, error) {
	var products []*model.ImportRequestProduct

	err := db.
		WithContext(ctx).
		Preload("SKU").
		Preload("Tote").
		Where(
			"request_id = ? AND tote_barcode = ?",
			requestID,
			toteBarcode,
		).
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	product *model.ImportRequestProduct,
) error {
	return db.
		WithContext(ctx).
		Model(&model.ImportRequestProduct{}).
		Where(
			"request_id = ? AND sku_barcode = ?",
			product.RequestID,
			product.SKUBarcode,
		).
		Updates(product).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
	skuBarcode string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.ImportRequestProduct{},
			"request_id = ? AND sku_barcode = ?",
			requestID,
			skuBarcode,
		).
		Error
}
