package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/sku"

	"github.com/gin-gonic/gin"
)

var _ = dto.SKUResponse{}

type SKUHandler struct {
	service *sku.Service
}

func NewSKUHandler(service *sku.Service) *SKUHandler {
	return &SKUHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create SKU
// @Description  Create a new SKU
// @Tags         skus
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSKURequest true "Create SKU request"
// @Success      201 {object} dto.SKUResponse
// @Failure      400 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /skus [post]
func (h *SKUHandler) Create(c *gin.Context) {
	var req dto.CreateSKURequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	sku, err := h.service.Create(
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
		"Create SKU successfully",
		sku,
	)
}

// FindByID godoc
// @Summary      Get SKU by barcode
// @Description  Get an SKU by its barcode
// @Tags         skus
// @Produce      json
// @Param        sku_barcode path string true "SKU barcode"
// @Success      200 {object} dto.SKUResponse
// @Failure      404 {object} gin.H
// @Router       /skus/{sku_barcode} [get]
func (h *SKUHandler) FindByID(c *gin.Context) {
	skuBarcode := c.Param("sku_barcode")

	sku, err := h.service.FindByID(
		c.Request.Context(),
		skuBarcode,
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
		"Get SKU successfully",
		sku,
	)
}

// FindAll godoc
// @Summary      Get all SKUs
// @Description  Get all SKUs
// @Tags         skus
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.SKUResponse
// @Failure      500 {object} gin.H
// @Router       /skus [get]
func (h *SKUHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	skus, err := h.service.FindAll(
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
		"Get SKUs successfully",
		skus,
		pagination,
	)
}

// Update godoc
// @Summary      Update SKU
// @Description  Update an existing SKU
// @Tags         skus
// @Accept       json
// @Produce      json
// @Param        sku_barcode path string true "SKU barcode"
// @Param        request body dto.UpdateSKURequest true "Update SKU request"
// @Success      200 {object} dto.SKUResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /skus/{sku_barcode} [put]
func (h *SKUHandler) Update(c *gin.Context) {
	skuBarcode := c.Param("sku_barcode")

	var req dto.UpdateSKURequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	sku, err := h.service.Update(
		c.Request.Context(),
		skuBarcode,
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
		"Update SKU successfully",
		sku,
	)
}

// Delete godoc
// @Summary      Delete SKU
// @Description  Delete an SKU by its barcode
// @Tags         skus
// @Produce      json
// @Param        sku_barcode path string true "SKU barcode"
// @Success      204
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /skus/{sku_barcode} [delete]
func (h *SKUHandler) Delete(c *gin.Context) {
	skuBarcode := c.Param("sku_barcode")

	err := h.service.Delete(
		c.Request.Context(),
		skuBarcode,
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
