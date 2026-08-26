package role

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/repository/role"
	"supermarket-backend/internal/response"
)

type Service struct {
	repository *role.Repository
}

func NewService(repository *role.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindByID(
	ctx context.Context,
	roleID string,
) (*dto.RoleResponse, error) {
	role, err := s.repository.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return dto.FromRoleModelToResponse(role), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.RoleResponse, error) {
	roles, err := s.repository.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.RoleResponse, 0, len(roles))

	for _, role := range roles {
		responses = append(
			responses,
			dto.FromRoleModelToResponse(role),
		)
	}

	return responses, nil
}
