package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	importRequestTote "supermarket-backend/internal/service/import_request_tote"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ = dto.ImportRequestToteResponse{}

type ImportRequestToteHandler struct {
	service *importRequestTote.Service
}

func NewImportRequestToteHandler(
	service *importRequestTote.Service,
) *ImportRequestToteHandler {
	return &ImportRequestToteHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create import request tote
// @Description  Add a tote to an import request
// @Tags         import-request-totes
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        tote_barcode path string true "Tote barcode"
// @Success      201 {object} dto.ImportRequestToteResponse
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/totes/{tote_barcode} [post]
func (h *ImportRequestToteHandler) Create(c *gin.Context) {
	requestID := c.Param("request_id")
	toteBarcode := c.Param("tote_barcode")

	tote, err := h.service.Create(
		c.Request.Context(),
		requestID,
		toteBarcode,
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

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Import request cannot be modified in its current status",
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
		http.StatusCreated,
		"Create import request tote successfully",
		tote,
	)
}

// FindAllByRequestID godoc
// @Summary      Get all totes of an import request
// @Description  Get all totes belonging to an import request
// @Tags         import-request-totes
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Success      200 {array} dto.ImportRequestToteResponse
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/totes [get]
func (h *ImportRequestToteHandler) FindAllByRequestID(c *gin.Context) {
	requestID := c.Param("request_id")

	totes, err := h.service.FindAllByRequestID(
		c.Request.Context(),
		requestID,
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
		http.StatusOK,
		"Get import request totes successfully",
		totes,
	)
}

// FindByBarcode godoc
// @Summary      Get import request tote by barcode
// @Description  Get a tote belonging to an import request
// @Tags         import-request-totes
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        tote_barcode path string true "Tote barcode"
// @Success      200 {object} dto.ImportRequestToteResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/totes/{tote_barcode} [get]
func (h *ImportRequestToteHandler) FindByBarcode(c *gin.Context) {
	requestID := c.Param("request_id")
	toteBarcode := c.Param("tote_barcode")

	tote, err := h.service.FindByID(
		c.Request.Context(),
		requestID,
		toteBarcode,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Import request tote not found",
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
		"Get import request tote successfully",
		tote,
	)
}

// Delete godoc
// @Summary      Delete import request tote
// @Description  Delete a tote from an import request
// @Tags         import-request-totes
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        tote_barcode path string true "Tote barcode"
// @Success      204
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/totes/{tote_barcode} [delete]
func (h *ImportRequestToteHandler) Delete(c *gin.Context) {
	requestID := c.Param("request_id")
	toteBarcode := c.Param("tote_barcode")

	err := h.service.Delete(
		c.Request.Context(),
		requestID,
		toteBarcode,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Import request or tote not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Import request cannot be modified in its current status",
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

	c.Status(http.StatusNoContent)
}
