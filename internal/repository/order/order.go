package order

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/response"

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
	order *model.Order,
) error {
	return db.
		WithContext(ctx).
		Create(order).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
) (*model.Order, error) {
	var order model.Order

	err := db.
		WithContext(ctx).
		Preload("Branch").
		Preload("Cashier").
		Preload("Items").
		Preload("Items.SKU").
		Preload("Items.SKU.Unit").
		Preload("Payments").
		First(&order, "order_id = ?", orderID).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
	query *dto.FindAllOrdersQuery,
	pagination *response.Pagination,
) ([]*model.Order, error) {
	var orders []*model.Order
	var total int64

	queryBuilder := db.
		WithContext(ctx).
		Model(&model.Order{})

	// Filter
	if query != nil && query.BranchID != nil {
		queryBuilder = queryBuilder.Where(
			"branch_id = ?",
			*query.BranchID,
		)
	}

	if query != nil && query.Status != nil {
		queryBuilder = queryBuilder.Where(
			"status = ?",
			*query.Status,
		)
	}

	// Count total records after applying filter
	if err := queryBuilder.
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

	// Get data
	if err := queryBuilder.
		Preload("Branch").
		Preload("Cashier").
		Limit(pagination.PerPage).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).
		Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	order *model.Order,
) error {
	return db.
		WithContext(ctx).
		Model(&model.Order{}).
		Where("order_id = ?", order.OrderID).
		Updates(order).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Order{},
			"order_id = ?",
			orderID,
		).
		Error
}
