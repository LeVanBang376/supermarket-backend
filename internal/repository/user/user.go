package user

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
	user *model.User,
) error {
	return db.
		WithContext(ctx).
		Create(user).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
) (*model.User, error) {
	var user model.User

	err := db.
		WithContext(ctx).
		Preload("Branch").
		Preload("Role").
		Preload("Position").
		First(&user, "user_id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByUsername(
	ctx context.Context,
	db *gorm.DB,
	username string,
) (*model.User, error) {
	var user model.User

	err := db.
		WithContext(ctx).
		Preload("Branch").
		Preload("Role").
		Preload("Position").
		First(&user, "username = ?", username).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
	db *gorm.DB,
) ([]*model.User, error) {
	var users []*model.User

	err := db.
		WithContext(ctx).
		Preload("Branch").
		Preload("Role").
		Preload("Position").
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	user *model.User,
) error {
	return db.
		WithContext(ctx).
		Save(user).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.User{},
			"user_id = ?",
			userID,
		).
		Error
}
