package auth

import "time"

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountResponse struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Role        Role       `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	BannedUntil *time.Time `json:"banned_until"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

type LoginResponse struct {
	AccessToken string           `json:"access_token"`
	Account     *AccountResponse `json:"account"`
}
