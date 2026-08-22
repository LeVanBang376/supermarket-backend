package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"supermarket-backend/internal/response"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				// Log panic + stack trace
				log.Printf(
					"panic recovered: %v\n%s",
					err,
					debug.Stack(),
				)

				response.NonDataJSON(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
