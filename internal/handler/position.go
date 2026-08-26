package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/position"

	"github.com/gin-gonic/gin"
)

var _ = dto.PositionResponse{}

type PositionHandler struct {
	service *position.Service
}

func NewPositionHandler(service *position.Service) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}

// FindByID godoc
// @Summary      Get position by ID
// @Description  Get a position by its ID
// @Tags         positions
// @Produce      json
// @Param        position_id path string true "Position ID"
// @Success      200 {object} dto.PositionResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /positions/{position_id} [get]
func (h *PositionHandler) FindByID(c *gin.Context) {
	positionID := c.Param("position_id")

	position, err := h.service.FindByID(
		c.Request.Context(),
		positionID,
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
		"Get position successfully",
		position,
	)
}

// FindAll godoc
// @Summary      Get all positions
// @Description  Get all positions
// @Tags         positions
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.PositionResponse
// @Failure      500 {object} gin.H
// @Router       /positions [get]
func (h *PositionHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	positions, err := h.service.FindAll(
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
		"Get positions successfully",
		positions,
		pagination,
	)
}
