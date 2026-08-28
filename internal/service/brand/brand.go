package brand

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/repository/brand"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *brand.Repository
}

func NewService(
	db *gorm.DB,
	repository *brand.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateBrandRequest,
) (*dto.BrandResponse, error) {
	var brand *model.Brand

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		brandID, err := s.repository.GetNextBrandID(ctx, tx)
		if err != nil {
			return err
		}

		brand = &model.Brand{
			BrandID:   brandID,
			BrandName: req.BrandName,
		}

		if err := s.repository.Create(ctx, tx, brand); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	brandID string,
) (*dto.BrandResponse, error) {
	brand, err := s.repository.FindByID(
		ctx,
		s.db,
		brandID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.BrandResponse, error) {
	brands, err := s.repository.FindAll(
		ctx,
		s.db,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.BrandResponse,
		0,
		len(brands),
	)

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
	brand, err := s.repository.FindByID(
		ctx,
		s.db,
		brandID,
	)
	if err != nil {
		return nil, err
	}

	if req.BrandName != nil {
		brand.BrandName = *req.BrandName
	}

	err = s.repository.Update(
		ctx,
		s.db,
		brand,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromBrandModelToResponse(brand), nil
}

func (s *Service) Delete(
	ctx context.Context,
	brandID string,
) error {
	_, err := s.repository.FindByID(
		ctx,
		s.db,
		brandID,
	)
	if err != nil {
		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		brandID,
	)
}
