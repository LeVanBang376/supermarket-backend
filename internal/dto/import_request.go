package dto

import (
	"time"

	"supermarket-backend/internal/model"

	"github.com/google/uuid"
)

type FindAllImportRequestsQuery struct {
	Status *string `form:"status"`
}

type CreateImportRequestRequest struct {
	BranchID           string                       `json:"branch_id" binding:"required"`
	ExpectedDeliveryAt time.Time                    `json:"expected_delivery_at" binding:"required"`
	Products           []CreateImportRequestProduct `json:"products" binding:"required,min=1"`
}

type CreateImportRequestProduct struct {
	SKUBarcode string `json:"sku_barcode" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,gt=0"`
}

func (r CreateImportRequestRequest) ToImportRequestModel(createdBy uuid.UUID, requestID string) *model.ImportRequest {
	return &model.ImportRequest{
		RequestID:          requestID,
		BranchID:           r.BranchID,
		CreatedBy:          createdBy,
		ExpectedDeliveryAt: r.ExpectedDeliveryAt,
		Status:             model.ImportRequestStatusDraft,
	}
}

func (r CreateImportRequestProduct) ToImportRequestProductModel(
	requestID string,
) *model.ImportRequestProduct {
	return &model.ImportRequestProduct{
		RequestID:  requestID,
		SKUBarcode: r.SKUBarcode,
		Quantity:   r.Quantity,
	}
}

type UpdateImportRequestRequest struct {
	BranchID                *string                      `json:"branch_id"`
	ExpectedDeliveryAt      *time.Time                   `json:"expected_delivery_at"`
	DeliveryLicensePlate    *string                      `json:"delivery_license_plate"`
	ReceivedBy              *uuid.UUID                   `json:"received_by"`
	Products                []UpdateImportRequestProduct `json:"products"`
	ProductsBarcodeToDelete []string                     `json:"products_barcode_to_delete"`
}

type UpdateImportRequestProduct struct {
	SKUBarcode       string  `json:"sku_barcode" binding:"required"`
	Quantity         *int    `json:"quantity,omitempty" binding:"omitempty,gt=0"`
	ToteBarcode      *string `json:"tote_barcode,omitempty"`
	LoadedQuantity   *int    `json:"loaded_quantity,omitempty" binding:"omitempty,gte=0"`
	ReceivedQuantity *int    `json:"received_quantity,omitempty" binding:"omitempty,gte=0"`
}

func (r UpdateImportRequestProduct) ToImportRequestProductModel(
	requestID string,
) *model.ImportRequestProduct {
	product := &model.ImportRequestProduct{
		RequestID:  requestID,
		SKUBarcode: r.SKUBarcode,
	}

	if r.Quantity != nil {
		product.Quantity = *r.Quantity
	}

	if r.ToteBarcode != nil {
		product.ToteBarcode = r.ToteBarcode
	}

	if r.LoadedQuantity != nil {
		product.LoadedQuantity = *r.LoadedQuantity
	}

	if r.ReceivedQuantity != nil {
		product.ReceivedQuantity = *r.ReceivedQuantity
	}

	return product
}

type UpdateImportRequestStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type ImportRequestProductResponse struct {
	SKUBarcode       string `json:"sku_barcode"`
	SKUName          string `json:"sku_name"`
	Quantity         int    `json:"quantity"`
	ToteBarcode      string `json:"tote_barcode,omitempty"`
	LoadedQuantity   int    `json:"loaded_quantity"`
	ReceivedQuantity int    `json:"received_quantity"`
}

type ImportRequestToteResponse struct {
	ToteBarcode string `json:"tote_barcode"`
}

type ImportRequestResponse struct {
	RequestID            string                         `json:"request_id"`
	Branch               *BranchResponse                `json:"branch"`
	CreatedBy            uuid.UUID                      `json:"created_by"`
	Creator              *UserResponse                  `json:"creator"`
	ExpectedDeliveryAt   time.Time                      `json:"expected_delivery_at"`
	DeliveryLicensePlate string                         `json:"delivery_license_plate"`
	Status               string                         `json:"status"`
	ReceivedBy           *uuid.UUID                     `json:"received_by,omitempty"`
	Receiver             *UserResponse                  `json:"receiver,omitempty"`
	CompleteAt           *time.Time                     `json:"complete_at,omitempty"`
	Totes                []ImportRequestToteResponse    `json:"totes"`
	Products             []ImportRequestProductResponse `json:"products"`
	CreatedAt            time.Time                      `json:"created_at"`
	UpdatedAt            time.Time                      `json:"updated_at"`
}

func FromImportRequestModelToResponse(
	request *model.ImportRequest,
) *ImportRequestResponse {
	return &ImportRequestResponse{
		RequestID:            request.RequestID,
		CreatedBy:            request.CreatedBy,
		ExpectedDeliveryAt:   request.ExpectedDeliveryAt,
		DeliveryLicensePlate: request.DeliveryLicensePlate,
		Status:               request.Status,
		ReceivedBy:           request.ReceivedBy,
		CompleteAt:           request.CompleteAt,
		CreatedAt:            request.CreatedAt,
		UpdatedAt:            request.UpdatedAt,
		Branch:               FromBranchModelToResponse(&request.Branch),
		Creator:              FromUserModelToResponse(&request.Creator),
		Receiver:             FromUserModelToResponse(request.Receiver),
	}
}

func FromImportRequestProductModelToResponse(
	product *model.ImportRequestProduct,
) *ImportRequestProductResponse {
	response := &ImportRequestProductResponse{
		SKUBarcode:       product.SKUBarcode,
		Quantity:         product.Quantity,
		LoadedQuantity:   product.LoadedQuantity,
		ReceivedQuantity: product.ReceivedQuantity,
	}

	if product.ToteBarcode != nil {
		response.ToteBarcode = *product.ToteBarcode
	}

	if product.SKU != nil {
		response.SKUName = product.SKU.SKUName
	}

	return response
}

func FromImportRequestToteModelToResponse(
	tote *model.ImportRequestTote,
) *ImportRequestToteResponse {
	return &ImportRequestToteResponse{
		ToteBarcode: tote.ToteBarcode,
	}
}
