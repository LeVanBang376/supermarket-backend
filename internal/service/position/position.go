package position

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/repository/position"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *position.Repository
}

func NewService(
	db *gorm.DB,
	repository *position.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) FindByID(
	ctx context.Context,
	positionID string,
) (*dto.PositionResponse, error) {
	position, err := s.repository.FindByID(
		ctx,
		s.db,
		positionID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromPositionModelToResponse(position), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.PositionResponse, error) {
	positions, err := s.repository.FindAll(
		ctx,
		s.db,
		pagination,
	)
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
