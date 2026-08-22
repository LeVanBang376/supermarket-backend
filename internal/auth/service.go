package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrEmailTaken         = errors.New("email already exists")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) error {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	account := &Account{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = s.repo.CreateAccount(ctx, account)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.ConstraintName, "username") {
					return ErrUsernameTaken
				}

				if strings.Contains(pgErr.ConstraintName, "email") {
					return ErrEmailTaken
				}
			}
		}

		return err
	}

	return nil
}

func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
) (*LoginResponse, error) {
	account, err := s.repo.FindByUsername(
		ctx,
		req.Username,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := GenerateToken(
		account.ID,
		account.Username,
		account.Role,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		Account: &AccountResponse{
			ID:          account.ID,
			Username:    account.Username,
			Email:       account.Email,
			Role:        account.Role,
			CreatedAt:   account.CreatedAt,
			BannedUntil: account.BannedUntil,
			LastLoginAt: account.LastLoginAt,
		},
	}, nil
}
