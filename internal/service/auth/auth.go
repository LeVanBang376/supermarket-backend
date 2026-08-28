package auth

import (
	"context"
	"errors"

	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/internal/dto"
	userRepository "supermarket-backend/internal/repository/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	userRepo   *userRepository.Repository
	jwtService *jwt.JWTService
}

func NewService(
	db *gorm.DB,
	userRepo *userRepository.Repository,
	jwtService *jwt.JWTService,
) *Service {
	return &Service{
		db:         db,
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *Service) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(
		ctx,
		s.db,
		req.Username,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != "ACTIVE" {
		return nil, errors.New("user is inactive")
	}

	accessToken, err := s.jwtService.GenerateToken(
		user.UserID,
		user.Username,
		user.RoleID,
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
		User:        *dto.FromUserModelToResponse(user),
	}, nil
}

func (s *Service) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(
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

	return dto.FromUserModelToResponse(user), nil
}
