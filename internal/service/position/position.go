package position

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/repository/position"
	"supermarket-backend/internal/response"
)

type Service struct {
	repository *position.Repository
}

func NewService(repository *position.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindByID(
	ctx context.Context,
	positionID string,
) (*dto.PositionResponse, error) {
	position, err := s.repository.FindByID(ctx, positionID)
	if err != nil {
		return nil, err
	}

	return dto.FromPositionModelToResponse(position), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.PositionResponse, error) {
	positions, err := s.repository.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.PositionResponse,
		0,
		len(positions),
	)

	for _, position := range positions {
		responses = append(
			responses,
			dto.FromPositionModelToResponse(position),
		)
	}

	return responses, nil
}
