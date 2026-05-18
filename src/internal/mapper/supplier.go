package mapper

import (
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/dto"
)

func SupplierModelToResponse(s *models.Supplier, addr dto.AddressResponse) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:          s.ID,
		Name:        s.Name,
		PhoneNumber: s.PhoneNumber,
		Address:     addr,
	}
}

func SupplierCreateRequestToModel(req dto.SupplierCreateRequest, addressID uuid.UUID) models.Supplier {
	return models.Supplier{
		ID:        uuid.New(),
		Name:      req.Name,
		AddressID: addressID,
		PhoneNumber: req.PhoneNumber,
	}
}