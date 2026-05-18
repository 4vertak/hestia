package dto

import "github.com/google/uuid"

type ImageResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id,omitempty"`
}

type ImageUploadRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Image     []byte    `json:"image" binding:"required"`
}