package unit

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/repository/unit"
	"supermarket-backend/internal/response"
)

type Service struct {
	repository *unit.Repository
}

func NewService(repository *unit.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateUnitRequest,
) (*dto.UnitResponse, error) {
	unitID, err := s.repository.GetNextUnitID(ctx)
	if err != nil {
		return nil, err
	}

	unit := &model.Unit{
		UnitID:   unitID,
		UnitName: req.UnitName,
	}

	err = s.repository.Create(ctx, unit)
	if err != nil {
		return nil, err
	}

	return dto.FromUnitModelToResponse(unit), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	unitID string,
) (*dto.UnitResponse, error) {
	unit, err := s.repository.FindByID(ctx, unitID)
	if err != nil {
		return nil, err
	}

	return dto.FromUnitModelToResponse(unit), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.UnitResponse, error) {
	units, err := s.repository.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UnitResponse, 0, len(units))

	for _, unit := range units {
		responses = append(
			responses,
			dto.FromUnitModelToResponse(unit),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	unitID string,
	req *dto.UpdateUnitRequest,
) (*dto.UnitResponse, error) {
	unit, err := s.repository.FindByID(ctx, unitID)
	if err != nil {
		return nil, err
	}

	if req.UnitName != nil {
		unit.UnitName = *req.UnitName
	}

	err = s.repository.Update(ctx, unit)
	if err != nil {
		return nil, err
	}

	return dto.FromUnitModelToResponse(unit), nil
}

func (s *Service) Delete(
	ctx context.Context,
	unitID string,
) error {
	_, err := s.repository.FindByID(ctx, unitID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, unitID)
}
