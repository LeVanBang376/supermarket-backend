package import_request

import (
	"context"
	"errors"
	"fmt"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	importRequestRepository "supermarket-backend/internal/repository/import_request"
	importRequestProductRepository "supermarket-backend/internal/repository/import_request_product"
	stockRepository "supermarket-backend/internal/repository/stock"
	"supermarket-backend/internal/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrLoadedQuantityMismatch = errors.New(
		"Số lượng hàng đã xếp phải bằng số lượng yêu cầu",
	)
	ErrProductToteBarcodeRequired = errors.New(
		"Tất cả sản phẩm phải được gán tote barcode",
	)
)

type Service struct {
	db                *gorm.DB
	repository        *importRequestRepository.Repository
	productRepository *importRequestProductRepository.Repository
	stockRepository   *stockRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *importRequestRepository.Repository,
	productRepository *importRequestProductRepository.Repository,
	stockRepository *stockRepository.Repository,
) *Service {
	return &Service{
		db:                db,
		repository:        repository,
		productRepository: productRepository,
		stockRepository:   stockRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateImportRequestRequest,
	createdBy uuid.UUID,
) (*dto.ImportRequestResponse, error) {
	var importRequest *model.ImportRequest

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		requestID, err := s.repository.GetNextRequestID(ctx, tx)
		if err != nil {
			return err
		}

		importRequest = req.ToImportRequestModel(
			createdBy,
			requestID,
		)

		if err := s.repository.Create(ctx, tx, importRequest); err != nil {
			return err
		}

		for _, productReq := range req.Products {
			product := productReq.ToImportRequestProductModel(requestID)

			if err := s.productRepository.Create(ctx, tx, product); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.FromImportRequestModelToResponse(importRequest), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	query dto.FindAllImportRequestsQuery,
	pagination *response.Pagination,
) ([]*dto.ImportRequestResponse, error) {
	requests, err := s.repository.FindAll(
		ctx,
		s.db,
		query,
		pagination,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.ImportRequestResponse,
		0,
		len(requests),
	)

	for _, request := range requests {
		responses = append(
			responses,
			dto.FromImportRequestModelToResponse(request),
		)
	}

	return responses, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	requestID string,
) (*dto.ImportRequestResponse, error) {
	importRequest, err := s.repository.FindByID(
		ctx,
		s.db,
		requestID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromImportRequestModelToResponse(importRequest), nil
}

func (s *Service) UpdateStatus(
	ctx context.Context,
	requestID string,
	req *dto.UpdateImportRequestStatusRequest,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find existing import request
		request, err := s.repository.FindByID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 2. Validate status transition
		switch request.Status {
		case model.ImportRequestStatusDraft:
			if req.Status != model.ImportRequestStatusRequired &&
				req.Status != model.ImportRequestStatusCancelled {
				return gorm.ErrInvalidData
			}

		case model.ImportRequestStatusRequired:
			if req.Status != model.ImportRequestStatusSupplierReceived &&
				req.Status != model.ImportRequestStatusCancelled {
				return gorm.ErrInvalidData
			}

		case model.ImportRequestStatusSupplierReceived:
			if req.Status != model.ImportRequestStatusDelivering {
				return gorm.ErrInvalidData
			}

			if err := s.validateSupplierReceivedProducts(
				ctx,
				tx,
				requestID,
			); err != nil {
				return err
			}

		case model.ImportRequestStatusDelivering:
			if req.Status != model.ImportRequestStatusRejected {
				return gorm.ErrInvalidData
			}

		case model.ImportRequestStatusCancelled,
			model.ImportRequestStatusRejected,
			model.ImportRequestStatusCompleted:
			return gorm.ErrInvalidData

		default:
			return gorm.ErrInvalidData
		}

		// 3. Update status
		request.Status = req.Status

		return s.repository.Update(
			ctx,
			tx,
			request,
		)
	})
}

func (s *Service) validateSupplierReceivedProducts(
	ctx context.Context,
	tx *gorm.DB,
	requestID string,
) error {
	products, err := s.productRepository.FindByRequestID(
		ctx,
		tx,
		requestID,
	)
	if err != nil {
		return err
	}

	for _, product := range products {
		// Loaded quantity must equal requested quantity.
		if product.LoadedQuantity != product.Quantity {
			return fmt.Errorf(
				"sản phẩm %s: số lượng đã xếp (%d) phải bằng số lượng yêu cầu (%d)",
				product.SKUBarcode,
				product.LoadedQuantity,
				product.Quantity,
			)
		}

		// Every product must have a tote barcode.
		if product.ToteBarcode == nil || *product.ToteBarcode == "" {
			return fmt.Errorf(
				"sản phẩm %s: chưa được gán tote barcode",
				product.SKUBarcode,
			)
		}
	}

	return nil
}

func (s *Service) Update(
	ctx context.Context,
	requestID string,
	req *dto.UpdateImportRequestRequest,
) (*dto.ImportRequestResponse, error) {
	var importRequest *model.ImportRequest

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find existing import request
		var err error
		importRequest, err = s.repository.FindByID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 2. Update import request fields
		switch importRequest.Status {
		case model.ImportRequestStatusDraft:
			if req.BranchID != nil {
				importRequest.BranchID = *req.BranchID
			}

			if req.ExpectedDeliveryAt != nil {
				importRequest.ExpectedDeliveryAt = *req.ExpectedDeliveryAt
			}

		case model.ImportRequestStatusCancelled,
			model.ImportRequestStatusRequired,
			model.ImportRequestStatusRejected,
			model.ImportRequestStatusCompleted:
			return gorm.ErrInvalidData

		case model.ImportRequestStatusSupplierReceived,
			model.ImportRequestStatusDelivering:
			// No import request fields can be updated.

		default:
			return gorm.ErrInvalidData
		}

		// 3. Update import request
		if err := s.repository.Update(
			ctx,
			tx,
			importRequest,
		); err != nil {
			return err
		}

		// 4. Update / create products
		for _, productReq := range req.Products {
			switch importRequest.Status {
			case model.ImportRequestStatusDraft:
				if err := s.updateOrCreateDraftProduct(
					ctx,
					tx,
					requestID,
					productReq,
				); err != nil {
					return err
				}

			case model.ImportRequestStatusSupplierReceived:
				if err := s.updateSupplierReceivedProduct(
					ctx,
					tx,
					requestID,
					productReq,
				); err != nil {
					return err
				}

			case model.ImportRequestStatusDelivering:
				if err := s.updateDeliveringProduct(
					ctx,
					tx,
					requestID,
					productReq,
				); err != nil {
					return err
				}

			case model.ImportRequestStatusRequired:
				// No product update.

			default:
				return gorm.ErrInvalidData
			}
		}

		// 5. Delete products
		// Product deletion is only allowed in DRAFT.
		if len(req.ProductsBarcodeToDelete) > 0 {
			if importRequest.Status != model.ImportRequestStatusDraft {
				return gorm.ErrInvalidData
			}

			for _, skuBarcode := range req.ProductsBarcodeToDelete {
				if err := s.productRepository.Delete(
					ctx,
					tx,
					requestID,
					skuBarcode,
				); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.FromImportRequestModelToResponse(
		importRequest,
	), nil
}

func (s *Service) updateOrCreateDraftProduct(
	ctx context.Context,
	tx *gorm.DB,
	requestID string,
	req dto.UpdateImportRequestProduct,
) error {
	product, err := s.productRepository.FindByID(
		ctx,
		tx,
		requestID,
		req.SKUBarcode,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			product := req.ToImportRequestProductModel(requestID)

			return s.productRepository.Create(
				ctx,
				tx,
				product,
			)
		}

		return err
	}

	// Existing product → update only quantity.
	if req.Quantity != nil {
		product.Quantity = *req.Quantity
	}

	return s.productRepository.Update(
		ctx,
		tx,
		product,
	)
}

func (s *Service) updateSupplierReceivedProduct(
	ctx context.Context,
	tx *gorm.DB,
	requestID string,
	req dto.UpdateImportRequestProduct,
) error {
	product, err := s.productRepository.FindByID(
		ctx,
		tx,
		requestID,
		req.SKUBarcode,
	)
	if err != nil {
		return err
	}

	if req.ToteBarcode != nil {
		product.ToteBarcode = req.ToteBarcode
	}

	if req.LoadedQuantity != nil {
		product.LoadedQuantity = *req.LoadedQuantity
	}

	return s.productRepository.Update(
		ctx,
		tx,
		product,
	)
}

func (s *Service) updateDeliveringProduct(
	ctx context.Context,
	tx *gorm.DB,
	requestID string,
	req dto.UpdateImportRequestProduct,
) error {
	product, err := s.productRepository.FindByID(
		ctx,
		tx,
		requestID,
		req.SKUBarcode,
	)
	if err != nil {
		return err
	}

	if req.ReceivedQuantity != nil {
		product.ReceivedQuantity = *req.ReceivedQuantity
	}

	return s.productRepository.Update(
		ctx,
		tx,
		product,
	)
}

func (s *Service) ConfirmImport(
	ctx context.Context,
	requestID string,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find existing import request
		request, err := s.repository.FindByID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 2. Only allow confirmation when request is DELIVERING
		if request.Status != model.ImportRequestStatusDelivering {
			return gorm.ErrInvalidData
		}

		// 3. Find all products
		products, err := s.productRepository.FindByRequestID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 4. Update stock
		for _, product := range products {
			// Validate received quantity
			if product.ReceivedQuantity < 0 ||
				product.ReceivedQuantity > product.LoadedQuantity {
				return gorm.ErrInvalidData
			}

			if product.ReceivedQuantity == 0 {
				continue
			}

			if err := s.stockRepository.IncreaseOrCreate(
				ctx,
				tx,
				request.BranchID,
				product.SKUBarcode,
				product.ReceivedQuantity,
			); err != nil {
				return err
			}
		}

		// 5. Update import request status
		request.Status = model.ImportRequestStatusCompleted

		return s.repository.Update(
			ctx,
			tx,
			request,
		)
	})
}
