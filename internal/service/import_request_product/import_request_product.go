package import_request_product

import (
	"context"

	"supermarket-backend/internal/dto"
	importRequestProductRepository "supermarket-backend/internal/repository/import_request_product"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *importRequestProductRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *importRequestProductRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) FindByRequestID(
	ctx context.Context,
	requestID string,
) ([]*dto.ImportRequestProductResponse, error) {
	products, err := s.repository.FindByRequestID(
		ctx,
		s.db,
		requestID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.ImportRequestProductResponse,
		0,
		len(products),
	)

	for _, product := range products {
		responses = append(
			responses,
			dto.FromImportRequestProductModelToResponse(product),
		)
	}

	return responses, nil
}

func (s *Service) FindByRequestAndToteBarcode(
	ctx context.Context,
	requestID string,
	toteBarcode string,
) ([]*dto.ImportRequestProductResponse, error) {
	products, err := s.repository.FindByRequestIDAndToteBarcode(
		ctx,
		s.db,
		requestID,
		toteBarcode,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.ImportRequestProductResponse,
		0,
		len(products),
	)

	for _, product := range products {
		responses = append(
			responses,
			dto.FromImportRequestProductModelToResponse(product),
		)
	}

	return responses, nil
}
