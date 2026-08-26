package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/branch"

	"github.com/gin-gonic/gin"
)

var _ = dto.BranchResponse{}

type BranchHandler struct {
	service *branch.Service
}

func NewBranchHandler(service *branch.Service) *BranchHandler {
	return &BranchHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create branch
// @Description  Create a new branch
// @Tags         branches
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateBranchRequest true "Create branch request"
// @Success      201 {object} dto.BranchResponse
// @Failure      400 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /branches [post]
func (h *BranchHandler) Create(c *gin.Context) {
	var req dto.CreateBranchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	branch, err := h.service.Create(
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
		"Create branch successfully",
		branch,
	)
}

// FindByID godoc
// @Summary      Get branch by ID
// @Description  Get a branch by its ID
// @Tags         branches
// @Produce      json
// @Param        branch_id path string true "Branch ID"
// @Success      200 {object} dto.BranchResponse
// @Failure      404 {object} gin.H
// @Router       /branches/{branch_id} [get]
func (h *BranchHandler) FindByID(c *gin.Context) {
	branchID := c.Param("branch_id")

	branch, err := h.service.FindByID(
		c.Request.Context(),
		branchID,
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
		branch,
	)
}

// FindAll godoc
// @Summary      Get all branches
// @Description  Get all branches
// @Tags         branches
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.BranchResponse
// @Failure      500 {object} gin.H
// @Router       /branches [get]
func (h *BranchHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	branches, err := h.service.FindAll(
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
		"Get branches successfully",
		branches,
		pagination,
	)
}

// Update godoc
// @Summary      Update branch
// @Description  Update an existing branch
// @Tags         branches
// @Accept       json
// @Produce      json
// @Param        branch_id path string true "Branch ID"
// @Param        request body dto.UpdateBranchRequest true "Update branch request"
// @Success      200 {object} dto.BranchResponse
// @Failure      400 {object} gin.H
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /branches/{branch_id} [put]
func (h *BranchHandler) Update(c *gin.Context) {
	branchID := c.Param("branch_id")

	var req dto.UpdateBranchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	branch, err := h.service.Update(
		c.Request.Context(),
		branchID,
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
		"Update branch successfully",
		branch,
	)
}

// Delete godoc
// @Summary      Delete branch
// @Description  Delete a branch by its ID
// @Tags         branches
// @Produce      json
// @Param        branch_id path string true "Branch ID"
// @Success      204
// @Failure      404 {object} gin.H
// @Failure      500 {object} gin.H
// @Router       /branches/{branch_id} [delete]
func (h *BranchHandler) Delete(c *gin.Context) {
	branchID := c.Param("branch_id")

	err := h.service.Delete(
		c.Request.Context(),
		branchID,
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
