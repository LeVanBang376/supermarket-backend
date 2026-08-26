package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/role"

	"github.com/gin-gonic/gin"
)

var _ = dto.RoleResponse{}

type RoleHandler struct {
	service *role.Service
}

func NewRoleHandler(service *role.Service) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

// FindByID godoc
// @Summary      Get role by ID
// @Description  Get a role by its ID
// @Tags         roles
// @Produce      json
// @Param        role_id path string true "Role ID"
// @Success      200 {object} dto.RoleResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /roles/{role_id} [get]
func (h *RoleHandler) FindByID(c *gin.Context) {
	roleID := c.Param("role_id")

	role, err := h.service.FindByID(
		c.Request.Context(),
		roleID,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get branch successfully",
		role,
	)
}

// FindAll godoc
// @Summary      Get all roles
// @Description  Get all roles
// @Tags         roles
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.RoleResponse
// @Failure      500 {object} gin.H
// @Router       /roles [get]
func (h *RoleHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	roles, err := h.service.FindAll(
		c.Request.Context(),
		pagination,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.PaginatedJSON(
		c.Writer,
		http.StatusOK,
		"Get roles successfully",
		roles,
		pagination,
	)
}
