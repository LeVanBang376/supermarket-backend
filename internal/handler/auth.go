package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/auth"

	"github.com/gin-gonic/gin"
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
// @Description  Authenticate a user and set the access token in an HttpOnly cookie
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

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			err.Error(),
		)
		return
	}

	c.SetCookie(
		"access_token",  // cookie name
		res.AccessToken, // JWT
		1200,            // MaxAge: 20 minutes
		"/",             // Path
		"",              // Domain
		false,           // Secure
		true,            // HttpOnly
	)

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Login successful",
		res,
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
