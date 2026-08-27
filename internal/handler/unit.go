package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/unit"

	"github.com/gin-gonic/gin"
)

var _ = dto.UnitResponse{}

type UnitHandler struct {
	service *unit.Service
}

func NewUnitHandler(service *unit.Service) *UnitHandler {
	return &UnitHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create unit
// @Description  Create a new unit
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUnitRequest true "Create unit request"
// @Success      201 {object} dto.UnitResponse
// @Failure      400 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /units [post]
func (h *UnitHandler) Create(c *gin.Context) {
	var req dto.CreateUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	unit, err := h.service.Create(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusCreated,
		"Create unit successfully",
		unit,
	)
}

// FindByID godoc
// @Summary      Get unit by ID
// @Description  Get a unit by its ID
// @Tags         units
// @Produce      json
// @Param        unit_id path string true "Unit ID"
// @Success      200 {object} dto.UnitResponse
// @Failure      404 {object} gin.H
// @Router       /units/{unit_id} [get]
func (h *UnitHandler) FindByID(c *gin.Context) {
	unitID := c.Param("unit_id")

	unit, err := h.service.FindByID(
		c.Request.Context(),
		unitID,
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
		"Get unit successfully",
		unit,
	)
}

// FindAll godoc
// @Summary      Get all units
// @Description  Get all units
// @Tags         units
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.UnitResponse
// @Failure      500 {object} gin.H
// @Router       /units [get]
func (h *UnitHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	units, err := h.service.FindAll(
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
		"Get units successfully",
		units,
		pagination,
	)
}

// Update godoc
// @Summary      Update unit
// @Description  Update an existing unit
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        unit_id path string true "Unit ID"
// @Param        request body dto.UpdateUnitRequest true "Update unit request"
// @Success      200 {object} dto.UnitResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /units/{unit_id} [put]
func (h *UnitHandler) Update(c *gin.Context) {
	unitID := c.Param("unit_id")

	var req dto.UpdateUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	unit, err := h.service.Update(
		c.Request.Context(),
		unitID,
		&req,
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
		"Update unit successfully",
		unit,
	)
}

// Delete godoc
// @Summary      Delete unit
// @Description  Delete a unit by its ID
// @Tags         units
// @Produce      json
// @Param        unit_id path string true "Unit ID"
// @Success      204
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /units/{unit_id} [delete]
func (h *UnitHandler) Delete(c *gin.Context) {
	unitID := c.Param("unit_id")

	err := h.service.Delete(
		c.Request.Context(),
		unitID,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	c.Status(http.StatusNoContent)
}
