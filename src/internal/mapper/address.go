package mapper

import (
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/models"
	"github.com/google/uuid"
)

func AddressModelToResponse(m *models.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:      m.ID,
		Country: m.Country,
		City:    m.City,
		Street:  m.Street,
	}
}

func AddressRequestToModel(id uuid.UUID, req dto.AddressRequest) models.Address {
	return models.Address{
		ID:      id,
		Country: req.Country,
		City:    req.City,
		Street:  req.Street,
	}
}