package stock

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/repository/stock"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *stock.Repository
}

func NewService(
	db *gorm.DB,
	repository *stock.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.StockResponse, error) {
	stocks, err := s.repository.FindAll(
		ctx,
		s.db,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.StockResponse,
		0,
		len(stocks),
	)

	for _, stock := range stocks {
		responses = append(
			responses,
			dto.FromStockModelToResponse(stock),
		)
	}

	return responses, nil
}
