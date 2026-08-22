package middleware

import (
	"context"
	"net/http"
	"strings"

	"supermarket-backend/internal/auth"
	"supermarket-backend/internal/response"
)

type contextKey string

const ContextUserKey contextKey = "user"

func Auth(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		authHeader := r.Header.Get(
			"Authorization",
		)

		if authHeader == "" {
			response.NonDataJSON(
				w,
				http.StatusUnauthorized,
				"Missing authorization header",
			)
			return
		}

		parts := strings.Split(
			authHeader,
			" ",
		)

		if len(parts) != 2 {
			response.NonDataJSON(
				w,
				http.StatusUnauthorized,
				"Invalid authorization header",
			)
			return
		}

		if parts[0] != "Bearer" {
			response.NonDataJSON(
				w,
				http.StatusUnauthorized,
				"Invalid auth scheme",
			)
			return
		}

		tokenString := parts[1]

		claims, err := auth.ParseToken(
			tokenString,
		)
		if err != nil {
			response.NonDataJSON(
				w,
				http.StatusUnauthorized,
				"Invalid token",
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			ContextUserKey,
			claims,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func GetClaims(
	ctx context.Context,
) *auth.Claims {
	claims, ok := ctx.Value(
		ContextUserKey,
	).(*auth.Claims)

	if !ok {
		return nil
	}

	return claims
}
