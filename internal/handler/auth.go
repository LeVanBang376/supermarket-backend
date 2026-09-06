package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/auth"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"

	accessTokenMaxAge  = 1200
	refreshTokenMaxAge = 7 * 24 * 60 * 60
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and set access and refresh tokens in HttpOnly cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest   true  "Login credentials"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	res, err := h.service.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			err.Error(),
		)
		return
	}

	// Access token
	c.SetCookie(
		accessTokenCookie,
		res.AccessToken,
		accessTokenMaxAge,
		"/",
		"",
		false,
		true,
	)

	// Refresh token
	c.SetCookie(
		refreshTokenCookie,
		res.RefreshToken,
		refreshTokenMaxAge,
		"/",
		"",
		false,
		true,
	)

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Login successful",
		&dto.LoginResponse{
			AccessToken: res.AccessToken,
			User:        *res.User,
		},
	)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Refresh the access and refresh tokens using the refresh token stored in an HttpOnly cookie
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  auth.RefreshResult
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Refresh token not found",
		)
		return
	}

	// Refresh tokens
	res, err := h.service.Refresh(
		c.Request.Context(),
		refreshToken,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			err.Error(),
		)
		return
	}

	// Set new access token
	c.SetCookie(
		accessTokenCookie,
		res.AccessToken,
		accessTokenMaxAge,
		"/",
		"",
		false,
		true,
	)

	// Set new refresh token
	c.SetCookie(
		refreshTokenCookie,
		res.RefreshToken,
		refreshTokenMaxAge,
		"/",
		"",
		false,
		true,
	)

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Token refreshed successfully",
	)
}

// Logout godoc
// @Summary      User logout
// @Description  Logout the current user by revoking the refresh token and clearing authentication cookies
// @Tags         Auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		// Không có refresh token thì vẫn clear cookies
		// và coi như logout thành công.
		h.clearAuthCookies(c)

		response.NonDataJSON(
			c.Writer,
			http.StatusOK,
			"Logout successful",
		)
		return
	}

	// Revoke refresh token session
	if err := h.service.Logout(
		c.Request.Context(),
		refreshToken,
	); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Logout failed",
		)
		return
	}

	// Clear authentication cookies
	h.clearAuthCookies(c)

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Logout successful",
	)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie(
		accessTokenCookie,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		refreshTokenCookie,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}

// Me godoc
// @Summary      Get current user
// @Description  Get information of the currently authenticated user
// @Tags         Auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)

	if claims == nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	user, err := h.service.GetUserByID(
		c.Request.Context(),
		claims.UserID,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get current user successful",
		user,
	)
}
