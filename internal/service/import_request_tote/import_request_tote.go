package import_request_tote

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	importRequestRepository "supermarket-backend/internal/repository/import_request"
	importRequestToteRepository "supermarket-backend/internal/repository/import_request_tote"

	"gorm.io/gorm"
)

type Service struct {
	db                *gorm.DB
	repository        *importRequestToteRepository.Repository
	importRequestRepo *importRequestRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *importRequestToteRepository.Repository,
	importRequestRepo *importRequestRepository.Repository,
) *Service {
	return &Service{
		db:                db,
		repository:        repository,
		importRequestRepo: importRequestRepo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	requestID string,
	toteBarcode string,
) (*dto.ImportRequestToteResponse, error) {
	var tote *model.ImportRequestTote

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find existing import request
		importRequest, err := s.importRequestRepo.FindByID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 2. Only allow adding tote when request is SUPPLIER_RECEIVED
		if importRequest.Status != model.ImportRequestStatusSupplierReceived {
			return gorm.ErrInvalidData
		}

		// 3. Create tote
		tote = &model.ImportRequestTote{
			RequestID:   requestID,
			ToteBarcode: toteBarcode,
		}

		if err := s.repository.Create(
			ctx,
			tx,
			tote,
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.FromImportRequestToteModelToResponse(tote), nil
}

func (s *Service) FindAllByRequestID(
	ctx context.Context,
	requestID string,
) ([]*dto.ImportRequestToteResponse, error) {
	totes, err := s.repository.FindByRequestID(
		ctx,
		s.db,
		requestID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.ImportRequestToteResponse,
		0,
		len(totes),
	)

	for _, tote := range totes {
		responses = append(
			responses,
			dto.FromImportRequestToteModelToResponse(tote),
		)
	}

	return responses, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	requestID string,
	toteBarcode string,
) (*dto.ImportRequestToteResponse, error) {
	tote, err := s.repository.FindByID(
		ctx,
		s.db,
		requestID,
		toteBarcode,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromImportRequestToteModelToResponse(tote), nil
}

func (s *Service) Delete(
	ctx context.Context,
	requestID string,
	toteBarcode string,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find existing import request
		importRequest, err := s.importRequestRepo.FindByID(
			ctx,
			tx,
			requestID,
		)
		if err != nil {
			return err
		}

		// 2. Only allow deleting tote when request is SUPPLIER_RECEIVED
		if importRequest.Status != model.ImportRequestStatusSupplierReceived {
			return gorm.ErrInvalidData
		}

		// 3. Delete tote
		return s.repository.Delete(
			ctx,
			tx,
			requestID,
			toteBarcode,
		)
	})
}
