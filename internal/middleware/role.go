package middleware

import (
	"net/http"

	"supermarket-backend/internal/auth"
	"supermarket-backend/internal/response"
)

func RequireRole(
	role auth.Role,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				claims := GetClaims(r.Context())
				if claims == nil {
					response.NonDataJSON(
						w,
						http.StatusUnauthorized,
						"Unauthorized",
					)
					return
				}

				if claims.Role != role {
					response.NonDataJSON(
						w,
						http.StatusForbidden,
						"You don't have permission to access this resource",
					)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
