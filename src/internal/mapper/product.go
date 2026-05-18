package mapper

import (
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/dto"
)

func ProductModelToResponse(p *models.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:             p.ID,
		Name:           p.Name,
		Category:       p.Category,
		Price:          p.Price,
		AvailableStock: p.AvailableStock,
		LastUpdateDate: p.LastUpdateDate,
		SupplierID:     p.SupplierID,
		ImageID:        p.ImageID,
	}
}

func ProductCreateRequestToModel(req dto.ProductCreateRequest, id uuid.UUID) models.Product {
	return models.Product{
		ID:             id,
		Name:           req.Name,
		Category:       req.Category,
		Price:          req.Price,
		AvailableStock: req.AvailableStock,
		SupplierID:     req.SupplierID,
		LastUpdateDate: nil, // будет проставлен сервисом
	}
}