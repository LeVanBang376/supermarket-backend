package role

import (
	"context"

	"supermarket-backend/internal/model"

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

func (r *Repository) Create(
	ctx context.Context,
	role *model.Role,
) error {
	return r.db.
		WithContext(ctx).
		Create(role).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	roleID string,
) (*model.Role, error) {
	var role model.Role

	err := r.db.
		WithContext(ctx).
		First(&role, "role_id = ?", roleID).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *Repository) FindAll(
	ctx context.Context,
) ([]*model.Role, error) {
	var roles []*model.Role

	err := r.db.
		WithContext(ctx).
		Find(&roles).
		Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) Update(
	ctx context.Context,
	role *model.Role,
) error {
	return r.db.
		WithContext(ctx).
		Model(&model.Role{}).
		Where("role_id = ?", role.RoleID).
		Updates(role).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	roleID string,
) error {
	return r.db.
		WithContext(ctx).
		Delete(
			&model.Role{},
			"role_id = ?",
			roleID,
		).
		Error
}
