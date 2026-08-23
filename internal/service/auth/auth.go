package auth

import (
	"context"
	"errors"

	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/internal/dto"
	userRepository "supermarket-backend/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	userRepo   *userRepository.Repository
	jwtService *jwt.JWTService
}

func NewService(
	userRepo *userRepository.Repository,
	jwtService *jwt.JWTService,
) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *Service) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
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
		User: dto.UserResponse{
			UserID:      user.UserID,
			Username:    user.Username,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			BranchID:    user.BranchID,
			RoleID:      user.RoleID,
			PositionID:  user.PositionID,
			Status:      user.Status,
		},
	}, nil
}
