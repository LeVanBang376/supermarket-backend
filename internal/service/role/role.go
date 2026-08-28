package role

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/repository/role"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *role.Repository
}

func NewService(
	db *gorm.DB,
	repository *role.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) FindByID(
	ctx context.Context,
	roleID string,
) (*dto.RoleResponse, error) {
	role, err := s.repository.FindByID(
		ctx,
		s.db,
		roleID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromRoleModelToResponse(role), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.RoleResponse, error) {
	roles, err := s.repository.FindAll(
		ctx,
		s.db,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.RoleResponse,
		0,
		len(roles),
	)

	for _, role := range roles {
		responses = append(
			responses,
			dto.FromRoleModelToResponse(role),
		)
	}

	return responses, nil
}
