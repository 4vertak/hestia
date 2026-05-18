package dto

import (
	"time"
	"github.com/google/uuid"
)

type ClientResponse struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Surname          string          `json:"surname"`
	Birthday         *time.Time      `json:"birthday,omitempty"`
	Gender           *string         `json:"gender,omitempty"`
	RegistrationDate time.Time       `json:"registration_date"`
	Address          AddressResponse `json:"address"`
}

type ClientCreateRequest struct {
	Name     string         `json:"name" binding:"required,max=100"`
	Surname  string         `json:"surname" binding:"required,max=100"`
	Birthday *time.Time     `json:"birthday"`                         // ISO 8601
	Gender   *string        `json:"gender" binding:"omitempty,oneof=M F O"`
	Address  AddressRequest `json:"address" binding:"required"`
}

type ClientUpdateAddressRequest struct {
	Address AddressRequest `json:"address" binding:"required"`
}

type ClientSearchParams struct {
	Name    string `form:"name" binding:"required_without=Surname"`
	Surname string `form:"surname" binding:"required_without=Name"`
}