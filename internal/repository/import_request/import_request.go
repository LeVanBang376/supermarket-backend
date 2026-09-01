package import_request

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

func (r *Repository) GetNextRequestID(
	ctx context.Context,
	db *gorm.DB,
) (string, error) {
	var nextID int64

	err := db.
		WithContext(ctx).
		Raw("SELECT nextval('import_request_id_seq')").
		Scan(&nextID).
		Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("IR%05d", nextID), nil
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	request *model.ImportRequest,
) error {
	return db.
		WithContext(ctx).
		Create(request).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	requestID string,
) (*model.ImportRequest, error) {
	var request model.ImportRequest

	err := db.
		WithContext(ctx).
		Preload("Branch").
		Preload("Creator").
		Preload("Receiver").
		First(&request, "request_id = ?", requestID).
		Error

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	pagination *response.Pagination,
) ([]*model.ImportRequest, error) {
	var requests []*model.ImportRequest

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.ImportRequest{}).
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
		Preload("Branch").
		Preload("Creator").
		Limit(pagination.PerPage).
		Offset(offset).
		Find(&requests).
		Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	request *model.ImportRequest,
) error {
	return db.
		WithContext(ctx).
		Model(&model.ImportRequest{}).
		Where("request_id = ?", request.RequestID).
		Updates(request).
		Error
}
