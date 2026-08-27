package dto

import "supermarket-backend/internal/model"

type CreateSKURequest struct {
	SKUBarcode    string `json:"sku_barcode" binding:"required,max=30"`
	SKUName       string `json:"sku_name" binding:"required,max=50"`
	BrandID       string `json:"brand_id" binding:"required,max=5"`
	UnitID        string `json:"unit_id" binding:"required,max=5"`
	ShelfLifeDays int    `json:"shelf_life_days" binding:"required,min=0"`
}

type UpdateSKURequest struct {
	SKUName       *string `json:"sku_name" binding:"omitempty,max=50"`
	BrandID       *string `json:"brand_id" binding:"omitempty,max=5"`
	UnitID        *string `json:"unit_id" binding:"omitempty,max=5"`
	ShelfLifeDays *int    `json:"shelf_life_days" binding:"omitempty,min=0"`
}

type SKUResponse struct {
	SKUBarcode    string `json:"sku_barcode"`
	SKUName       string `json:"sku_name"`
	BrandID       string `json:"brand_id"`
	UnitID        string `json:"unit_id"`
	ShelfLifeDays int    `json:"shelf_life_days"`
}

func FromSKUModelToResponse(sku *model.SKU) *SKUResponse {
	var s *SKUResponse = &SKUResponse{}

	s.SKUBarcode = sku.SKUBarcode
	s.SKUName = sku.SKUName
	s.BrandID = sku.BrandID
	s.UnitID = sku.UnitID
	s.ShelfLifeDays = sku.ShelfLifeDays

	return s
}
