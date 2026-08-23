package dto

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
