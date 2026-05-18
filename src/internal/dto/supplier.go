package dto

import "github.com/google/uuid"

type SupplierResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	PhoneNumber *string         `json:"phone_number,omitempty"`
	Address     AddressResponse `json:"address"`
}

type SupplierCreateRequest struct {
	Name        string         `json:"name" binding:"required,max=255"`
	PhoneNumber *string        `json:"phone_number" binding:"omitempty,max=20"`
	Address     AddressRequest `json:"address" binding:"required"`
}

type SupplierUpdateAddressRequest struct {
	Address AddressRequest `json:"address" binding:"required"`
}