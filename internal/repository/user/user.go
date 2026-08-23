package user

import (
	"context"

	"supermarket-backend/internal/model"

	"github.com/google/uuid"
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

func (r *Repository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).
		Create(user).
		Error
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).
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

func (r *Repository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.db.WithContext(ctx).
		First(&user, "username = ?", username).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	err := r.db.WithContext(ctx).
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

func (r *Repository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).
		Save(user).
		Error
}

func (r *Repository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.User{}, "user_id = ?", userID).
		Error
}
