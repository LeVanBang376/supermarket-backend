package sku

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	skuRepository "supermarket-backend/internal/repository/sku"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *skuRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *skuRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateSKURequest,
) (*dto.SKUResponse, error) {
	sku := &model.SKU{
		SKUBarcode:    req.SKUBarcode,
		SKUName:       req.SKUName,
		BrandID:       req.BrandID,
		UnitID:        req.UnitID,
		ShelfLifeDays: req.ShelfLifeDays,
	}

	err := s.repository.Create(
		ctx,
		s.db,
		sku,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromSKUModelToResponse(sku), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	skuBarcode string,
) (*dto.SKUResponse, error) {
	sku, err := s.repository.FindByID(
		ctx,
		s.db,
		skuBarcode,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromSKUModelToResponse(sku), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.SKUResponse, error) {
	skus, err := s.repository.FindAll(
		ctx,
		s.db,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.SKUResponse,
		0,
		len(skus),
	)

	for _, sku := range skus {
		responses = append(
			responses,
			dto.FromSKUModelToResponse(sku),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	skuBarcode string,
	req *dto.UpdateSKURequest,
) (*dto.SKUResponse, error) {
	sku, err := s.repository.FindByID(
		ctx,
		s.db,
		skuBarcode,
	)
	if err != nil {
		return nil, err
	}

	if req.SKUName != nil {
		sku.SKUName = *req.SKUName
	}

	if req.BrandID != nil {
		sku.BrandID = *req.BrandID
	}

	if req.UnitID != nil {
		sku.UnitID = *req.UnitID
	}

	if req.ShelfLifeDays != nil {
		sku.ShelfLifeDays = *req.ShelfLifeDays
	}

	err = s.repository.Update(
		ctx,
		s.db,
		sku,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromSKUModelToResponse(sku), nil
}

func (s *Service) Delete(
	ctx context.Context,
	skuBarcode string,
) error {
	_, err := s.repository.FindByID(
		ctx,
		s.db,
		skuBarcode,
	)
	if err != nil {
		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		skuBarcode,
	)
}
