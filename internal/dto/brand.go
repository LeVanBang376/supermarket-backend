package dto

import "supermarket-backend/internal/model"

type CreateBrandRequest struct {
	BrandName string `json:"brand_name" binding:"required,max=30"`
}

type UpdateBrandRequest struct {
	BrandName *string `json:"brand_name" binding:"omitempty,max=30"`
}

type BrandResponse struct {
	BrandID   string `json:"brand_id"`
	BrandName string `json:"brand_name"`
}

func FromBrandModelToResponse(brand *model.Brand) *BrandResponse {
	var b *BrandResponse = &BrandResponse{}

	b.BrandID = brand.BrandID
	b.BrandName = brand.BrandName

	return b
}
