package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"supermarket-backend/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register godoc
// @Summary Register account
// @Description Create new account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register payload"
// @Success 201 {object} response.NonDataResponse
// @Failure 400 {object} response.NonDataResponse
// @Router /register [post]
func (h *Handler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	err = h.service.Register(
		r.Context(),
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameTaken):
			response.NonDataJSON(
				w,
				http.StatusConflict,
				err.Error(),
			)

		case errors.Is(err, ErrEmailTaken):
			response.NonDataJSON(
				w,
				http.StatusConflict,
				err.Error(),
			)

		default:
			response.NonDataJSON(
				w,
				http.StatusInternalServerError,
				"internal server error",
			)
		}

		return
	}

	response.NonDataJSON(
		w,
		http.StatusCreated,
		"Account is created successfully!",
	)
}

// Login godoc
// @Summary Login account
// @Description Authenticate account and return account data
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} response.Response[LoginResponse]
// @Failure 400 {object} response.NonDataResponse
// @Failure 401 {object} response.NonDataResponse
// @Failure 500 {object} response.NonDataResponse
// @Router /login [post]
func (h *Handler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	account, err := h.service.Login(
		r.Context(),
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.NonDataJSON(
				w,
				http.StatusUnauthorized,
				err.Error(),
			)

		default:
			response.NonDataJSON(
				w,
				http.StatusInternalServerError,
				"internal server error",
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Login successful",
		account,
	)
}
