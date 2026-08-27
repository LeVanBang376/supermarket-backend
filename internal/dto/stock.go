package dto

import "supermarket-backend/internal/model"

type StockResponse struct {
	BranchID   string `json:"branch_id"`
	SKUBarcode string `json:"sku_barcode"`
	Quantity   int    `json:"quantity"`
}

func FromStockModelToResponse(stock *model.Stock) *StockResponse {
	var s *StockResponse = &StockResponse{}

	s.BranchID = stock.BranchID
	s.SKUBarcode = stock.SKUBarcode
	s.Quantity = stock.Quantity

	return s
}
