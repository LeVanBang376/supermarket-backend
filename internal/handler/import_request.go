package handler

import (
	"errors"
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/middleware"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/import_request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.ImportRequestResponse{}

type ImportRequestHandler struct {
	service *import_request.Service
}

func NewImportRequestHandler(
	service *import_request.Service,
) *ImportRequestHandler {
	return &ImportRequestHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create import request
// @Description  Create a new import request
// @Tags         import-requests
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateImportRequestRequest true "Create import request"
// @Success      201 {object} dto.ImportRequestResponse
// @Failure      400 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests [post]
func (h *ImportRequestHandler) Create(c *gin.Context) {
	var req dto.CreateImportRequestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	createdBy := middleware.GetUserID(c)

	if createdBy == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	importRequest, err := h.service.Create(
		c.Request.Context(),
		&req,
		createdBy,
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
		"Create import request successfully",
		importRequest,
	)
}

// FindAll godoc
// @Summary      Get all import requests
// @Description  Get all import requests
// @Tags         import-requests
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.ImportRequestResponse
// @Failure      500 {object} gin.H
// @Router       /import-requests [get]
func (h *ImportRequestHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	requests, err := h.service.FindAll(
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
		"Get import requests successfully",
		requests,
		pagination,
	)
}

// FindByID godoc
// @Summary      Get import request by ID
// @Description  Get an import request by ID
// @Tags         import-requests
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Success      200 {object} dto.ImportRequestResponse
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id} [get]
func (h *ImportRequestHandler) FindByID(c *gin.Context) {
	requestID := c.Param("request_id")

	importRequest, err := h.service.FindByID(
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
		"Get import request successfully",
		importRequest,
	)
}

// Update godoc
// @Summary      Update import request
// @Description  Update an import request and its products
// @Tags         import-requests
// @Accept       json
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        request body dto.UpdateImportRequestRequest true "Update import request"
// @Success      200 {object} dto.ImportRequestResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id} [patch]
func (h *ImportRequestHandler) Update(c *gin.Context) {
	requestID := c.Param("request_id")

	var req dto.UpdateImportRequestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	importRequest, err := h.service.Update(
		c.Request.Context(),
		requestID,
		&req,
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
				"Import request cannot be updated in its current status",
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
		"Update import request successfully",
		importRequest,
	)
}

// UpdateStatus godoc
// @Summary      Update import request status
// @Description  Update the status of an import request
// @Tags         import-requests
// @Accept       json
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Param        request body dto.UpdateImportRequestStatusRequest true "Update import request status"
// @Success      200 {object} dto.ImportRequestResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/status [patch]
func (h *ImportRequestHandler) UpdateStatus(c *gin.Context) {
	requestID := c.Param("request_id")

	var req dto.UpdateImportRequestStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	err := h.service.UpdateStatus(
		c.Request.Context(),
		requestID,
		&req,
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
				"Invalid import request status transition",
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

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Update import request status successfully",
	)
}

// ConfirmImport godoc
// @Summary      Confirm import request
// @Description  Confirm an import request and update stock quantities
// @Tags         import-requests
// @Produce      json
// @Param        request_id path string true "Import request ID"
// @Success      200
// @Failure      404 {object} gin.H
// @Failure      409 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /import-requests/{request_id}/confirm [post]
func (h *ImportRequestHandler) ConfirmImport(c *gin.Context) {
	requestID := c.Param("request_id")

	err := h.service.ConfirmImport(
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

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Import request cannot be confirmed in its current state",
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

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Import request confirmed successfully",
	)
}
