package dto

import "supermarket-backend/internal/model"

type PositionResponse struct {
	PositionID   string `json:"position_id"`
	PositionName string `json:"position_name"`
}

func FromPositionModelToResponse(
	position *model.Position,
) *PositionResponse {
	return &PositionResponse{
		PositionID:   position.PositionID,
		PositionName: position.PositionName,
	}
}
