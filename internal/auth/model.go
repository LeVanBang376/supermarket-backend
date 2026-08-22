package auth

import "time"

type Role string

const (
	RolePlayer Role = "player"
	RoleAdmin  Role = "admin"
)

type Account struct {
	ID           int64
	Username     string
	PasswordHash string
	Email        string
	Role         Role
	CreatedAt    time.Time
	BannedUntil  *time.Time
	LastLoginAt  *time.Time
}
