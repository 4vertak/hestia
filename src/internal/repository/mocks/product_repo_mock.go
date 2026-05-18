package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ProductRepositoryMock struct {
	CreateFn           func(ctx context.Context, product *models.Product) error
	DecreaseStockFn    func(ctx context.Context, id uuid.UUID, quantity int) error
	GetByIDFn          func(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindAllAvailableFn func(ctx context.Context) ([]models.Product, error)
	DeleteFn           func(ctx context.Context, id uuid.UUID) error
	UpdateLastUpdateFn func(ctx context.Context, id uuid.UUID) error
	UpdateImageIDFn    func(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *models.Product) error {
	return m.CreateFn(ctx, product)
}

func (m *ProductRepositoryMock) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return m.DecreaseStockFn(ctx, id, quantity)
}

func (m *ProductRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *ProductRepositoryMock) FindAllAvailable(ctx context.Context) ([]models.Product, error) {
	return m.FindAllAvailableFn(ctx)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

func (m *ProductRepositoryMock) UpdateLastUpdate(ctx context.Context, id uuid.UUID) error {
	return m.UpdateLastUpdateFn(ctx, id)
}

func (m *ProductRepositoryMock) UpdateImageID(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error {
	return m.UpdateImageIDFn(ctx, productID, imageID)
}