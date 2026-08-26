package branch

import (
	"context"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	"supermarket-backend/internal/repository/branch"
	"supermarket-backend/internal/response"
)

type Service struct {
	repository *branch.Repository
}

func NewService(repository *branch.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateBranchRequest,
) (*dto.BranchResponse, error) {
	branchID, err := s.repository.GetNextBranchID(ctx)
	if err != nil {
		return nil, err
	}

	branch := &model.Branch{
		BranchID:   branchID,
		BranchName: req.BranchName,
		Address:    req.Address,
	}

	err = s.repository.Create(ctx, branch)
	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	branchID string,
) (*dto.BranchResponse, error) {
	branch, err := s.repository.FindByID(ctx, branchID)
	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) FindAll(
	ctx context.Context,
	pagination *response.Pagination,
) ([]*dto.BranchResponse, error) {
	branches, err := s.repository.FindAll(ctx, pagination)
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
	branch, err := s.repository.FindByID(ctx, branchID)
	if err != nil {
		return nil, err
	}

	if req.BranchName != nil {
		branch.BranchName = *req.BranchName
	}

	if req.Address != nil {
		branch.Address = *req.Address
	}

	err = s.repository.Update(ctx, branch)
	if err != nil {
		return nil, err
	}

	return dto.FromBranchModelToResponse(branch), nil
}

func (s *Service) Delete(
	ctx context.Context,
	branchID string,
) error {
	_, err := s.repository.FindByID(ctx, branchID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, branchID)
}
