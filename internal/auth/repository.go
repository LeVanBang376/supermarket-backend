package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAccount(
	ctx context.Context,
	account *Account,
) error {
	query := `
		INSERT INTO accounts (
			username,
			password_hash,
			email
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		account.Username,
		account.PasswordHash,
		account.Email,
	)

	return err
}

func (r *Repository) FindByUsername(
	ctx context.Context,
	username string,
) (*Account, error) {
	query := `
		SELECT
			id,
			username,
			password_hash,
			email,
			role,
			created_at,
			banned_until,
			last_login_at
		FROM accounts
		WHERE username = $1
	`

	var account Account

	err := r.db.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&account.ID,
		&account.Username,
		&account.PasswordHash,
		&account.Email,
		&account.Role,
		&account.CreatedAt,
		&account.BannedUntil,
		&account.LastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
