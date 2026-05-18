package dto

import "github.com/google/uuid"

type AddressResponse struct {
	ID      uuid.UUID `json:"id"`
	Country string    `json:"country"`
	City    string    `json:"city"`
	Street  string    `json:"street"`
}

type AddressRequest struct {
	Country string `json:"country" binding:"required,max=100"`
	City    string `json:"city" binding:"required,max=100"`
	Street  string `json:"street" binding:"required,max=255"`
}