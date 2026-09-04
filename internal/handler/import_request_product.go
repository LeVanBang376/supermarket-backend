package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	importRequestProduct "supermarket-backend/internal/service/import_request_product"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ = dto.ImportRequestProductResponse{}

type ImportRequestProductHandler struct {
	service *importRequestProduct.Service
}

func NewImportRequestProductHandler(
	service *importRequestProduct.Service,
) *ImportRequestProductHandler {
	return &ImportRequestProductHandler{
		service: service,
	}
}

// FindByRequestID godoc
// @Summary      Get products of an import request
// @Description  Get all products belonging to an import request
// @Tags         import-request-products
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Success      200 {array} dto.ImportRequestProductResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/products [get]
func (h *ImportRequestProductHandler) FindByRequestID(c *gin.Context) {
	requestID := c.Param("request_id")

	products, err := h.service.FindByRequestID(
		c.Request.Context(),
		requestID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Import request not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get import request products successfully",
		products,
	)
}

// FindByRequestAndToteBarcode godoc
// @Summary      Get products of a tote
// @Description  Get all products belonging to a specific tote in an import request
// @Tags         import-request-products
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        tote_barcode path string true "Tote barcode"
// @Success      200 {array} dto.ImportRequestProductResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/totes/{tote_barcode}/products [get]
func (h *ImportRequestProductHandler) FindByRequestAndToteBarcode(c *gin.Context) {
	requestID := c.Param("request_id")
	toteBarcode := c.Param("tote_barcode")

	products, err := h.service.FindByRequestAndToteBarcode(
		c.Request.Context(),
		requestID,
		toteBarcode,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Import request products not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get import request tote products successfully",
		products,
	)
}
