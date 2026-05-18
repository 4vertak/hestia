package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type SupplierRepositoryMock struct {
	CreateFn        func(ctx context.Context, supplier *models.Supplier) error
	UpdateAddressFn func(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error
	DeleteFn        func(ctx context.Context, id uuid.UUID) error
	FindAllFn       func(ctx context.Context) ([]models.SupplierWithAddress, error)
	GetByIDFn       func(ctx context.Context, id uuid.UUID) (*models.SupplierWithAddress, error)
}

func (m *SupplierRepositoryMock) Create(ctx context.Context, supplier *models.Supplier) error {
	return m.CreateFn(ctx, supplier)
}

func (m *SupplierRepositoryMock) UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error {
	return m.UpdateAddressFn(ctx, id, newAddressID)
}

func (m *SupplierRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

func (m *SupplierRepositoryMock) FindAll(ctx context.Context) ([]models.SupplierWithAddress, error) {
	return m.FindAllFn(ctx)
}

func (m *SupplierRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.SupplierWithAddress, error) {
	return m.GetByIDFn(ctx, id)
}