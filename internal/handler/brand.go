package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/brand"

	"github.com/gin-gonic/gin"
)

var _ = dto.BrandResponse{}

type BrandHandler struct {
	service *brand.Service
}

func NewBrandHandler(service *brand.Service) *BrandHandler {
	return &BrandHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create brand
// @Description  Create a new brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateBrandRequest true "Create brand request"
// @Success      201 {object} dto.BrandResponse
// @Failure      400 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /brands [post]
func (h *BrandHandler) Create(c *gin.Context) {
	var req dto.CreateBrandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	brand, err := h.service.Create(
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
		"Create brand successfully",
		brand,
	)
}

// FindByID godoc
// @Summary      Get brand by ID
// @Description  Get a brand by its ID
// @Tags         brands
// @Produce      json
// @Param        brand_id path string true "Brand ID"
// @Success      200 {object} dto.BrandResponse
// @Failure      404 {object} gin.H
// @Router       /brands/{brand_id} [get]
func (h *BrandHandler) FindByID(c *gin.Context) {
	brandID := c.Param("brand_id")

	brand, err := h.service.FindByID(
		c.Request.Context(),
		brandID,
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
		"Get brand successfully",
		brand,
	)
}

// FindAll godoc
// @Summary      Get all brands
// @Description  Get all brands
// @Tags         brands
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.BrandResponse
// @Failure      500 {object} gin.H
// @Router       /brands [get]
func (h *BrandHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	brands, err := h.service.FindAll(
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
		"Get brands successfully",
		brands,
		pagination,
	)
}

// Update godoc
// @Summary      Update brand
// @Description  Update an existing brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param        brand_id path string true "Brand ID"
// @Param        request body dto.UpdateBrandRequest true "Update brand request"
// @Success      200 {object} dto.BrandResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /brands/{brand_id} [put]
func (h *BrandHandler) Update(c *gin.Context) {
	brandID := c.Param("brand_id")

	var req dto.UpdateBrandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	brand, err := h.service.Update(
		c.Request.Context(),
		brandID,
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
		"Update brand successfully",
		brand,
	)
}

// Delete godoc
// @Summary      Delete brand
// @Description  Delete a brand by its ID
// @Tags         brands
// @Produce      json
// @Param        brand_id path string true "Brand ID"
// @Success      204
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /brands/{brand_id} [delete]
func (h *BrandHandler) Delete(c *gin.Context) {
	brandID := c.Param("brand_id")

	err := h.service.Delete(
		c.Request.Context(),
		brandID,
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
