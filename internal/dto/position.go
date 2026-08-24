package dto

import "supermarket-backend/internal/model"

type CreatePositionRequest struct {
	PositionName string `json:"position_name" binding:"required"`
}

type UpdatePositionRequest struct {
	PositionName *string `json:"position_name"`
}

type PositionResponse struct {
	PositionID   string `json:"position_id"`
	PositionName string `json:"position_name"`
}

func FromPositionModelToResponse(position *model.Position) *PositionResponse {
	var p *PositionResponse = &PositionResponse{}

	p.PositionID = position.PositionID
	p.PositionName = position.PositionName

	return p
}
