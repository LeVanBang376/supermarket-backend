package payment

import (
	"context"
	"time"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	orderRepository "supermarket-backend/internal/repository/order"
	paymentRepository "supermarket-backend/internal/repository/payment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db              *gorm.DB
	repository      *paymentRepository.Repository
	orderRepository *orderRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *paymentRepository.Repository,
	orderRepository *orderRepository.Repository,
) *Service {
	return &Service{
		db:              db,
		repository:      repository,
		orderRepository: orderRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	orderID uuid.UUID,
	req *dto.CreatePaymentRequest,
) (*dto.PaymentResponse, error) {
	var payment *model.Payment

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find order
		order, err := s.orderRepository.FindByID(
			ctx,
			tx,
			orderID,
		)
		if err != nil {
			return err
		}

		// 2. Only OPEN order can be paid
		if order.Status != model.OrderStatusOpen {
			return gorm.ErrInvalidData
		}

		// 3. Payment amount must match order total
		if req.Amount != order.TotalAmount {
			return gorm.ErrInvalidData
		}

		payment = &model.Payment{
			PaymentID:      uuid.New(),
			OrderID:        orderID,
			Method:         req.Method,
			Amount:         req.Amount,
			Status:         model.PaymentStatusSuccess,
			TransactionRef: req.TransactionRef,
		}

		now := time.Now()
		payment.PaidAt = &now

		// 4. Create payment
		if err := s.repository.Create(
			ctx,
			tx,
			payment,
		); err != nil {
			return err
		}

		// 5. Complete order
		order.Status = model.OrderStatusCompleted

		return s.orderRepository.Update(
			ctx,
			tx,
			order,
		)
	})

	if err != nil {
		return nil, err
	}

	return dto.FromPaymentModelToResponse(payment), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	paymentID uuid.UUID,
) (*dto.PaymentResponse, error) {
	payment, err := s.repository.FindByID(
		ctx,
		s.db,
		paymentID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromPaymentModelToResponse(payment), nil
}

func (s *Service) FindByOrderID(
	ctx context.Context,
	orderID uuid.UUID,
) ([]*dto.PaymentResponse, error) {
	payments, err := s.repository.FindByOrderID(
		ctx,
		s.db,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.PaymentResponse,
		0,
		len(payments),
	)

	for _, payment := range payments {
		responses = append(
			responses,
			dto.FromPaymentModelToResponse(payment),
		)
	}

	return responses, nil
}
