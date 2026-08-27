package handler

import (
	"net/http"

	"supermarket-backend/internal/dto"
	"supermarket-backend/internal/response"
	"supermarket-backend/internal/service/stock"

	"github.com/gin-gonic/gin"
)

var _ = dto.StockResponse{}

type StockHandler struct {
	service *stock.Service
}

func NewStockHandler(service *stock.Service) *StockHandler {
	return &StockHandler{
		service: service,
	}
}

// FindAll godoc
// @Summary      Get all stocks
// @Description  Get all stocks
// @Tags         stocks
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of items per page" default(10)
// @Success      200 {array} dto.StockResponse
// @Failure      500 {object} gin.H
// @Router       /stocks [get]
func (h *StockHandler) FindAll(c *gin.Context) {
	pagination := response.NewPagination(c.Request)

	stocks, err := h.service.FindAll(
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
		"Get stocks successfully",
		stocks,
		pagination,
	)
}
