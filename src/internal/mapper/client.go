package mapper

import (
	"time"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/dto"
)

func ClientModelToResponse(c *models.Client, addr dto.AddressResponse) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               c.ID,
		Name:             c.Name,
		Surname:          c.Surname,
		Birthday:         c.Birthday,
		Gender:           c.Gender,
		RegistrationDate: c.RegistrationDate,
		Address:          addr,
	}
}

func ClientCreateRequestToModel(req dto.ClientCreateRequest, addressID uuid.UUID) models.Client {
	return models.Client{
		ID:               uuid.New(),
		Name:             req.Name,
		Surname:          req.Surname,
		Birthday:         req.Birthday,
		Gender:           req.Gender,
		RegistrationDate: time.Now(),
		AddressID:        addressID,
	}
}