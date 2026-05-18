package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/mapper"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository"
)

type ClientService struct {
	clientRepo  repository.ClientRepository
	addressRepo repository.AddressRepository
}

func NewClientService(cr repository.ClientRepository, ar repository.AddressRepository) *ClientService {
	return &ClientService{clientRepo: cr, addressRepo: ar}
}

func (s *ClientService) Create(ctx context.Context, req dto.ClientCreateRequest) (*dto.ClientResponse, error) {
	addrID := uuid.New()
	address := mapper.AddressRequestToModel(addrID, req.Address)
	if err := s.addressRepo.Create(ctx, &address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	client := models.Client{
		ID:               uuid.New(),
		Name:             req.Name,
		Surname:          req.Surname,
		Birthday:         req.Birthday,
		Gender:           req.Gender,
		RegistrationDate: time.Now(),
		AddressID:        addrID,
	}
	if err := s.clientRepo.Create(ctx, &client); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	addrResp := mapper.AddressModelToResponse(&address)
	resp := mapper.ClientModelToResponse(&client, addrResp)
	return &resp, nil
}

func (s *ClientService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.clientRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *ClientService) SearchByName(ctx context.Context, name, surname string) ([]dto.ClientResponse, error) {
	clients, err := s.clientRepo.FindByNameAndSurname(ctx, name, surname)
	if err != nil {
		return nil, err
	}
	var result []dto.ClientResponse
	for _, c := range clients {
		addr := dto.AddressResponse{
			ID:      c.AddressID,
			Country: c.Country,
			City:    c.City,
			Street:  c.Street,
		}
		client := models.Client{
			ID:               c.ID,
			Name:             c.Name,
			Surname:          c.Surname,
			Birthday:         c.Birthday,
			Gender:           c.Gender,
			RegistrationDate: c.RegistrationDate,
			AddressID:        c.AddressID,
		}
		result = append(result, mapper.ClientModelToResponse(&client, addr))
	}
	return result, nil
}

func (s *ClientService) GetAll(ctx context.Context, limit, offset int) ([]dto.ClientResponse, error) {
	clients, err := s.clientRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var result []dto.ClientResponse
	for _, c := range clients {
		addr := dto.AddressResponse{
			ID:      c.AddressID,
			Country: c.Country,
			City:    c.City,
			Street:  c.Street,
		}
		client := models.Client{
			ID:               c.ID,
			Name:             c.Name,
			Surname:          c.Surname,
			Birthday:         c.Birthday,
			Gender:           c.Gender,
			RegistrationDate: c.RegistrationDate,
			AddressID:        c.AddressID,
		}
		result = append(result, mapper.ClientModelToResponse(&client, addr))
	}
	return result, nil
}

func (s *ClientService) UpdateAddress(ctx context.Context, id uuid.UUID, req dto.ClientUpdateAddressRequest) (*dto.ClientResponse, error) {
	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil || client == nil {
		return nil, ErrNotFound
	}

	oldAddrID := client.AddressID

	newAddrID := uuid.New()
	newAddr := mapper.AddressRequestToModel(newAddrID, req.Address)
	if err := s.addressRepo.Create(ctx, &newAddr); err != nil {
		return nil, err
	}

	err = s.clientRepo.UpdateAddress(ctx, id, newAddrID)
	if err != nil {
		return nil, err
	}

	// старый адрес удаляем и игнорируем ошибки
	_ = s.addressRepo.Delete(ctx, oldAddrID)

	addrResp := mapper.AddressModelToResponse(&newAddr)
	resp := mapper.ClientModelToResponse(client, addrResp)
	return &resp, nil
}