package order

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	orderRepository "supermarket-backend/internal/repository/order"
	"supermarket-backend/internal/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *orderRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *orderRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateOrderRequest,
	cashierID uuid.UUID,
) (*dto.OrderResponse, error) {
	order := &model.Order{
		OrderID:        uuid.New(),
		BranchID:       req.BranchID,
		CashierID:      cashierID,
		Status:         model.OrderStatusOpen,
		Subtotal:       0,
		DiscountAmount: 0,
		TotalAmount:    0,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		order,
	); err != nil {
		return nil, err
	}

	return dto.FromOrderModelToResponse(order), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	query *dto.FindAllOrdersQuery,
	pagination *response.Pagination,
) ([]*dto.OrderResponse, error) {
	orders, err := s.repository.FindAll(
		ctx,
		s.db,
		query,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.OrderResponse,
		0,
		len(orders),
	)

	for _, order := range orders {
		responses = append(
			responses,
			dto.FromOrderModelToResponse(order),
		)
	}

	return responses, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	orderID uuid.UUID,
) (*dto.OrderResponse, error) {
	order, err := s.repository.FindByID(
		ctx,
		s.db,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromOrderModelToResponse(order), nil
}

func (s *Service) Cancel(
	ctx context.Context,
	orderID uuid.UUID,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.repository.FindByID(
			ctx,
			tx,
			orderID,
		)
		if err != nil {
			return err
		}

		if order.Status != model.OrderStatusOpen {
			return gorm.ErrInvalidData
		}

		order.Status = model.OrderStatusCancelled

		return s.repository.Update(
			ctx,
			tx,
			order,
		)
	})
}
