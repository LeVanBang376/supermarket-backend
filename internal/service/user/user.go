package user

import (
	"context"
	"errors"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	userRepository "supermarket-backend/internal/repository/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *userRepository.Repository
}

func NewService(
	db *gorm.DB,
	repo *userRepository.Repository,
) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req dto.CreateUserRequest,
) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := req.ToUserModel(string(passwordHash))

	if err := s.repo.Create(
		ctx,
		s.db,
		user,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) FindByID(
	ctx context.Context,
	userID uuid.UUID,
) (*model.User, error) {
	user, err := s.repo.FindByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return user, nil
}

func (s *Service) FindAll(
	ctx context.Context,
) ([]*model.User, error) {
	return s.repo.FindAll(
		ctx,
		s.db,
	)
}

func (s *Service) Update(
	ctx context.Context,
	userID uuid.UUID,
	req dto.UpdateUserRequest,
) (*model.User, error) {
	user, err := s.repo.FindByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	if req.BranchID != nil {
		user.BranchID = *req.BranchID
	}

	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}

	if req.PositionID != nil {
		user.PositionID = *req.PositionID
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.Update(
		ctx,
		s.db,
		user,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Delete(
	ctx context.Context,
	userID uuid.UUID,
) error {
	if _, err := s.FindByID(
		ctx,
		userID,
	); err != nil {
		return err
	}

	return s.repo.Delete(
		ctx,
		s.db,
		userID,
	)
}
