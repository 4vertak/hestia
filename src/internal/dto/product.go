package dto

import (
	"time"
	"github.com/google/uuid"
)

type ProductResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Category       string     `json:"category"`
	Price          float64    `json:"price"`
	AvailableStock int        `json:"available_stock"`
	LastUpdateDate *time.Time `json:"last_update_date,omitempty"`
	SupplierID     uuid.UUID  `json:"supplier_id"`
	ImageID        *uuid.UUID `json:"image_id,omitempty"`
}

type ProductCreateRequest struct {
	Name           string    `json:"name" binding:"required,max=255"`
	Category       string    `json:"category" binding:"required,max=100"`
	Price          float64   `json:"price" binding:"required,gt=0"`
	AvailableStock int       `json:"available_stock" binding:"min=0"`
	SupplierID     uuid.UUID `json:"supplier_id" binding:"required"`
}

type ProductDecreaseRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}