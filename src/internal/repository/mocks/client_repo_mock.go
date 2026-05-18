package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ClientRepositoryMock struct {
	CreateFn                 func(ctx context.Context, client *models.Client) error
	DeleteFn                 func(ctx context.Context, id uuid.UUID) error
	FindByNameAndSurnameFn   func(ctx context.Context, name, surname string) ([]models.ClientWithAddress, error)
	FindAllFn                func(ctx context.Context, limit, offset int) ([]models.ClientWithAddress, error)
	UpdateAddressFn          func(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error
	GetByIDFn                func(ctx context.Context, id uuid.UUID) (*models.Client, error)
}

func (m *ClientRepositoryMock) Create(ctx context.Context, client *models.Client) error {
	return m.CreateFn(ctx, client)
}

func (m *ClientRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

func (m *ClientRepositoryMock) FindByNameAndSurname(ctx context.Context, name, surname string) ([]models.ClientWithAddress, error) {
	return m.FindByNameAndSurnameFn(ctx, name, surname)
}

func (m *ClientRepositoryMock) FindAll(ctx context.Context, limit, offset int) ([]models.ClientWithAddress, error) {
	return m.FindAllFn(ctx, limit, offset)
}

func (m *ClientRepositoryMock) UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error {
	return m.UpdateAddressFn(ctx, id, newAddressID)
}

func (m *ClientRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return m.GetByIDFn(ctx, id)
}