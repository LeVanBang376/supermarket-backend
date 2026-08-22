package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		h.Login(w, r)
	})
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		h.Register(w, r)
	})
}
