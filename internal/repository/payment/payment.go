package payment

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
	payment *model.Payment,
) error {
	return db.
		WithContext(ctx).
		Create(payment).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	paymentID uuid.UUID,
) (*model.Payment, error) {
	var payment model.Payment

	err := db.
		WithContext(ctx).
		First(
			&payment,
			"payment_id = ?",
			paymentID,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *Repository) FindByOrderID(
	ctx context.Context,
	db *gorm.DB,
	orderID uuid.UUID,
) ([]*model.Payment, error) {
	var payments []*model.Payment

	err := db.
		WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).
		Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	payment *model.Payment,
) error {
	return db.
		WithContext(ctx).
		Model(&model.Payment{}).
		Where("payment_id = ?", payment.PaymentID).
		Updates(payment).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	paymentID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Payment{},
			"payment_id = ?",
			paymentID,
		).
		Error
}
