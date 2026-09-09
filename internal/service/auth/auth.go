package auth

import (
	"context"
	"errors"
	"time"

	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/infrastructure/token"
	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/model"
	userRepository "supermarket-backend/internal/repository/user"
	userSessionRepository "supermarket-backend/internal/repository/user_session"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const refreshTokenDuration = 7 * 24 * time.Hour

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *dto.UserResponse
}

type Service struct {
	db              *gorm.DB
	userRepo        *userRepository.Repository
	userSessionRepo *userSessionRepository.Repository
	jwtService      *jwt.JWTService
}

func NewService(
	db *gorm.DB,
	userRepo *userRepository.Repository,
	userSessionRepo *userSessionRepository.Repository,
	jwtService *jwt.JWTService,
) *Service {
	return &Service{
		db:              db,
		userRepo:        userRepo,
		userSessionRepo: userSessionRepo,
		jwtService:      jwtService,
	}
}

func (s *Service) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*LoginResult, error) {
	user, err := s.userRepo.FindByUsername(
		ctx,
		s.db,
		req.Username,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("username không tồn tại")
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("Mật khẩu sai")
	}

	if user.Status != "ACTIVE" {
		return nil, errors.New("user is inactive")
	}

	accessToken, err := s.jwtService.GenerateToken(
		user.UserID,
		user.Username,
		user.BranchID,
		user.RoleID,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenHash := token.HashRefreshToken(refreshToken)

	session := &model.UserSession{
		ID:               uuid.New(),
		UserID:           user.UserID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(refreshTokenDuration),
	}

	if err := s.userSessionRepo.Create(
		ctx,
		s.db,
		session,
	); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.FromUserModelToResponse(user),
	}, nil
}

func (s *Service) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	refreshTokenHash := token.HashRefreshToken(refreshToken)

	session, err := s.userSessionRepo.FindByRefreshTokenHash(
		ctx,
		s.db,
		refreshTokenHash,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Token không tồn tại thì coi như đã logout.
			return nil
		}

		return err
	}

	if session.RevokedAt != nil {
		// Session đã bị revoke rồi.
		return nil
	}

	return s.userSessionRepo.Revoke(
		ctx,
		s.db,
		session.ID,
	)
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

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*RefreshResult, error) {
	// 1. Hash refresh token
	refreshTokenHash := token.HashRefreshToken(refreshToken)

	// 2. Find session
	session, err := s.userSessionRepo.FindByRefreshTokenHash(
		ctx,
		s.db,
		refreshTokenHash,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}

		return nil, err
	}

	// 3. Check revoked
	if session.RevokedAt != nil {
		return nil, errors.New("refresh token has been revoked")
	}

	// 4. Check expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	// 5. Find user
	user, err := s.userRepo.FindByID(
		ctx,
		s.db,
		session.UserID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	// 6. Check user status
	if user.Status != "ACTIVE" {
		return nil, errors.New("user is inactive")
	}

	// 7. Generate new access token
	accessToken, err := s.jwtService.GenerateToken(
		user.UserID,
		user.Username,
		user.BranchID,
		user.RoleID,
	)
	if err != nil {
		return nil, err
	}

	// 8. Generate new refresh token
	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshTokenHash := token.HashRefreshToken(
		newRefreshToken,
	)

	// 9. Create new session
	newSession := &model.UserSession{
		ID:               uuid.New(),
		UserID:           user.UserID,
		RefreshTokenHash: newRefreshTokenHash,
		ExpiresAt:        time.Now().Add(refreshTokenDuration),
	}

	// 10. Rotate refresh token
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Revoke old session
		if err := s.userSessionRepo.Revoke(
			ctx,
			tx,
			session.ID,
		); err != nil {
			return err
		}

		// Create new session
		if err := s.userSessionRepo.Create(
			ctx,
			tx,
			newSession,
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
