package import_request_tote

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
	tote *model.ImportRequestTote,
) error {
	return db.
		WithContext(ctx).
		Create(tote).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
	toteBarcode string,
) (*model.ImportRequestTote, error) {
	var tote model.ImportRequestTote

	err := db.
		WithContext(ctx).
		First(
			&tote,
			"request_id = ? AND tote_barcode = ?",
			requestID,
			toteBarcode,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &tote, nil
}

func (r *Repository) FindByRequestID(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
) ([]*model.ImportRequestTote, error) {
	var totes []*model.ImportRequestTote

	err := db.
		WithContext(ctx).
		Where("request_id = ?", requestID).
		Find(&totes).
		Error

	if err != nil {
		return nil, err
	}

	return totes, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	tote *model.ImportRequestTote,
) error {
	return db.
		WithContext(ctx).
		Model(&model.ImportRequestTote{}).
		Where(
			"request_id = ? AND tote_barcode = ?",
			tote.RequestID,
			tote.ToteBarcode,
		).
		Updates(tote).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
	toteBarcode string,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.ImportRequestTote{},
			"request_id = ? AND tote_barcode = ?",
			requestID,
			toteBarcode,
		).
		Error
}
