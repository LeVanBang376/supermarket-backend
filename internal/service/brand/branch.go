package brand

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/repository/brand"
	"supermarket-backend/internal/response"
)

type Service struct {
	repository *brand.Repository
}

func NewService(repository *brand.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateBrandRequest,
) (*dto.BrandResponse, error) {
	brandID, err := s.repository.GetNextBrandID(ctx)
	if err != nil {
		return nil, err
	}

	brand := &model.Brand{
		BrandID:   brandID,
		BrandName: req.BrandName,
	}

	err = s.repository.Create(ctx, brand)
	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	brandID string,
) (*dto.BrandResponse, error) {
	brand, err := s.repository.FindByID(ctx, brandID)
	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.BrandResponse, error) {
	brands, err := s.repository.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.BrandResponse, 0, len(brands))

	for _, brand := range brands {
		responses = append(
			responses,
			dto.FromBrandModelToResponse(brand),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	brandID string,
	req *dto.UpdateBrandRequest,
) (*dto.BrandResponse, error) {
	brand, err := s.repository.FindByID(ctx, brandID)
	if err != nil {
		return nil, err
	}

	if req.BrandName != nil {
		brand.BrandName = *req.BrandName
	}

	err = s.repository.Update(ctx, brand)
	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) Delete(
	ctx context.Context,
	brandID string,
) error {
	_, err := s.repository.FindByID(ctx, brandID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, brandID)
}
