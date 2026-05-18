package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type AddressRepositoryMock struct {
	CreateFn  func(ctx context.Context, addr *models.Address) error
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*models.Address, error)
	UpdateFn  func(ctx context.Context, addr *models.Address) error
	DeleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *AddressRepositoryMock) Create(ctx context.Context, addr *models.Address) error {
	return m.CreateFn(ctx, addr)
}

func (m *AddressRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *AddressRepositoryMock) Update(ctx context.Context, addr *models.Address) error {
	return m.UpdateFn(ctx, addr)
}

func (m *AddressRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}