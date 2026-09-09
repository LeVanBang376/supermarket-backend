package order_item

import (
	"context"
	"fmt"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	orderRepository "supermarket-backend/internal/repository/order"
	orderItemRepository "supermarket-backend/internal/repository/order_item"
	skuRepository "supermarket-backend/internal/repository/sku"
	"supermarket-backend/internal/ws"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Service struct {
	db              *gorm.DB
	hub             *ws.WebsocketHub
	repository      *orderItemRepository.Repository
	orderRepository *orderRepository.Repository
	skuRepository   *skuRepository.Repository
}

func NewService(
	db *gorm.DB,
	hub *ws.WebsocketHub,
	repository *orderItemRepository.Repository,
	orderRepository *orderRepository.Repository,
	skuRepository *skuRepository.Repository,
) *Service {
	return &Service{
		db:              db,
		hub:             hub,
		repository:      repository,
		orderRepository: orderRepository,
		skuRepository:   skuRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	actorID uuid.UUID,
	orderID uuid.UUID,
	req *dto.AddOrderItemRequest,
) (*dto.OrderItemResponse, error) {
	var item *model.OrderItem

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

		// 2. Only OPEN order can be modified
		if order.Status != model.OrderStatusOpen {
			return gorm.ErrInvalidData
		}

		// 3. Check SKU exists
		sku, err := s.skuRepository.FindByID(
			ctx,
			tx,
			req.SKUBarcode,
		)
		if err != nil {
			return err
		}

		// 4. Check whether item already exists
		_, err = s.repository.FindByID(
			ctx,
			tx,
			orderID,
			req.SKUBarcode,
		)

		if err == nil {
			// SKU already exists in order
			return gorm.ErrDuplicatedKey
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		// 5. Create new order item
		item = &model.OrderItem{
			OrderID:    orderID,
			SKUBarcode: sku.SKUBarcode,
			Quantity:   req.Quantity,
			UnitPrice:  sku.UnitPrice,
			Subtotal:   req.Quantity * sku.UnitPrice,
		}

		if err := s.repository.Create(
			ctx,
			tx,
			item,
		); err != nil {
			return err
		}

		// 6. Recalculate order totals
		if err := s.recalculateOrder(
			ctx,
			tx,
			order,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf(
		`{"type":"ORDER_UPDATED","order_id":"%s"}`,
		orderID,
	)

	if err := s.hub.Broadcast(
		ctx,
		message,
		actorID,
		ws.WsCustomerDisplayType,
	); err != nil {
		return nil, err
	}

	return dto.FromOrderItemModelToResponse(item), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	orderID uuid.UUID,
	skuBarcode string,
) (*dto.OrderItemResponse, error) {
	item, err := s.repository.FindByID(
		ctx,
		s.db,
		orderID,
		skuBarcode,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromOrderItemModelToResponse(item), nil
}

func (s *Service) FindByOrderID(
	ctx context.Context,
	orderID uuid.UUID,
) ([]*dto.OrderItemResponse, error) {
	items, err := s.repository.FindByOrderID(
		ctx,
		s.db,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.OrderItemResponse,
		0,
		len(items),
	)

	for _, item := range items {
		responses = append(
			responses,
			dto.FromOrderItemModelToResponse(item),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	actorID uuid.UUID,
	orderID uuid.UUID,
	skuBarcode string,
	req *dto.UpdateOrderItemRequest,
) (*dto.OrderItemResponse, error) {
	var item *model.OrderItem

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepository.FindByID(
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

		item, err = s.repository.FindByID(
			ctx,
			tx,
			orderID,
			skuBarcode,
		)
		if err != nil {
			return err
		}

		if req.Quantity != nil {
			item.Quantity = *req.Quantity
		}

		item.Subtotal = item.Quantity * item.UnitPrice

		if err := s.repository.Update(
			ctx,
			tx,
			item,
		); err != nil {
			return err
		}

		return s.recalculateOrder(
			ctx,
			tx,
			order,
		)
	})

	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf(
		`{"type":"ORDER_UPDATED","order_id":"%s"}`,
		orderID,
	)

	if err := s.hub.Broadcast(
		ctx,
		message,
		actorID,
		ws.WsCustomerDisplayType,
	); err != nil {
		return nil, err
	}

	return dto.FromOrderItemModelToResponse(item), nil
}

func (s *Service) Delete(
	ctx context.Context,
	actorID uuid.UUID,
	orderID uuid.UUID,
	skuBarcode string,
) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.orderRepository.FindByID(
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

		if err := s.repository.Delete(
			ctx,
			tx,
			orderID,
			skuBarcode,
		); err != nil {
			return err
		}

		return s.recalculateOrder(
			ctx,
			tx,
			order,
		)
	})

	if err != nil {
		return err
	}

	message := fmt.Sprintf(
		`{"type":"ORDER_UPDATED","order_id":"%s"}`,
		orderID,
	)

	if err := s.hub.Broadcast(
		ctx,
		message,
		actorID,
		ws.WsCustomerDisplayType,
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) recalculateOrder(
	ctx context.Context,
	tx *gorm.DB,
	order *model.Order,
) error {
	items, err := s.repository.FindByOrderID(
		ctx,
		tx,
		order.OrderID,
	)
	if err != nil {
		return err
	}

	var subtotal float64

	for _, item := range items {
		subtotal += item.Subtotal
	}

	order.Subtotal = subtotal
	order.TotalAmount = subtotal - order.DiscountAmount

	return s.orderRepository.Update(
		ctx,
		tx,
		order,
	)
}
