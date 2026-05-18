package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/mapper"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository"
)

type SupplierService struct {
	supplierRepo repository.SupplierRepository
	addressRepo  repository.AddressRepository
}

func NewSupplierService(sr repository.SupplierRepository, ar repository.AddressRepository) *SupplierService {
	return &SupplierService{supplierRepo: sr, addressRepo: ar}
}

func (s *SupplierService) Create(ctx context.Context, req dto.SupplierCreateRequest) (*dto.SupplierResponse, error) {
	addrID := uuid.New()
	addr := mapper.AddressRequestToModel(addrID, req.Address)
	if err := s.addressRepo.Create(ctx, &addr); err != nil {
		return nil, err
	}

	supplier := models.Supplier{
		ID:          uuid.New(),
		Name:        req.Name,
		AddressID:   addrID,
		PhoneNumber: req.PhoneNumber,
	}
	if err := s.supplierRepo.Create(ctx, &supplier); err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	addrResp := mapper.AddressModelToResponse(&addr)
	resp := mapper.SupplierModelToResponse(&supplier, addrResp)
	return &resp, nil
}

func (s *SupplierService) UpdateAddress(ctx context.Context, id uuid.UUID, req dto.SupplierUpdateAddressRequest) (*dto.SupplierResponse, error) {
	supplierWithAddr, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil || supplierWithAddr == nil {
		return nil, ErrNotFound
	}

	oldAddrID := supplierWithAddr.AddressID

	newAddrID := uuid.New()
	newAddr := mapper.AddressRequestToModel(newAddrID, req.Address)
	if err := s.addressRepo.Create(ctx, &newAddr); err != nil {
		return nil, err
	}

	err = s.supplierRepo.UpdateAddress(ctx, id, newAddrID)
	if err != nil {
		return nil, err
	}

	_ = s.addressRepo.Delete(ctx, oldAddrID)

	// ответ по старым данным поставщика, но с новым адресом
	supplier := models.Supplier{
		ID:          supplierWithAddr.ID,
		Name:        supplierWithAddr.Name,
		AddressID:   newAddrID,
		PhoneNumber: supplierWithAddr.PhoneNumber,
	}
	addrResp := mapper.AddressModelToResponse(&newAddr)
	resp := mapper.SupplierModelToResponse(&supplier, addrResp)
	return &resp, nil
}

func (s *SupplierService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.supplierRepo.Delete(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *SupplierService) GetAll(ctx context.Context) ([]dto.SupplierResponse, error) {
	suppliers, err := s.supplierRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var result []dto.SupplierResponse
	for _, srow := range suppliers {
		addr := dto.AddressResponse{
			ID:      srow.AddressID,
			Country: srow.Country,
			City:    srow.City,
			Street:  srow.Street,
		}
		supplier := models.Supplier{
			ID:          srow.ID,
			Name:        srow.Name,
			AddressID:   srow.AddressID,
			PhoneNumber: srow.PhoneNumber,
		}
		result = append(result, mapper.SupplierModelToResponse(&supplier, addr))
	}
	return result, nil
}

func (s *SupplierService) GetByID(ctx context.Context, id uuid.UUID) (*dto.SupplierResponse, error) {
	srow, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil || srow == nil {
		return nil, ErrNotFound
	}
	addr := dto.AddressResponse{
		ID:      srow.AddressID,
		Country: srow.Country,
		City:    srow.City,
		Street:  srow.Street,
	}
	supplier := models.Supplier{
		ID:          srow.ID,
		Name:        srow.Name,
		AddressID:   srow.AddressID,
		PhoneNumber: srow.PhoneNumber,
	}
	resp := mapper.SupplierModelToResponse(&supplier, addr)
	return &resp, nil
}