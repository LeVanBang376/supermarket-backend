package branch

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/repository/branch"
	"supermarket-backend/internal/response"

	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *branch.Repository
}

func NewService(db *gorm.DB, repository *branch.Repository) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateBranchRequest,
) (*dto.BranchResponse, error) {
	var branch *model.Branch

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		branchID, err := s.repository.GetNextBranchID(ctx, tx)
		if err != nil {
			return err
		}

		branch = &model.Branch{
			BranchID:   branchID,
			BranchName: req.BranchName,
			Address:    req.Address,
		}

		if err := s.repository.Create(ctx, tx, branch); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	branchID string,
) (*dto.BranchResponse, error) {
	branch, err := s.repository.FindByID(ctx, s.db, branchID)
	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.BranchResponse, error) {
	branches, err := s.repository.FindAll(ctx, s.db, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.BranchResponse, 0, len(branches))

	for _, branch := range branches {
		responses = append(
			responses,
			dto.FromBranchModelToResponse(branch),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	branchID string,
	req *dto.UpdateBranchRequest,
) (*dto.BranchResponse, error) {
	branch, err := s.repository.FindByID(ctx, s.db, branchID)
	if err != nil {
		return nil, err
	}

	if req.BranchName != nil {
		branch.BranchName = *req.BranchName
	}

	if req.Address != nil {
		branch.Address = *req.Address
	}

	err = s.repository.Update(ctx, s.db, branch)
	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) Delete(
	ctx context.Context,
	branchID string,
) error {
	_, err := s.repository.FindByID(ctx, s.db, branchID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, s.db, branchID)
}
