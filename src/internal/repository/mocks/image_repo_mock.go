package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ImageRepositoryMock struct {
	CreateFn         func(ctx context.Context, img *models.Image) error
	ReplaceFn        func(ctx context.Context, id uuid.UUID, newData []byte) error
	DeleteFn         func(ctx context.Context, id uuid.UUID) error
	GetByIDFn        func(ctx context.Context, id uuid.UUID) (*models.Image, error)
	GetByProductIDFn func(ctx context.Context, productID uuid.UUID) (*models.Image, error)
}

func (m *ImageRepositoryMock) Create(ctx context.Context, img *models.Image) error {
	return m.CreateFn(ctx, img)
}

func (m *ImageRepositoryMock) Replace(ctx context.Context, id uuid.UUID, newData []byte) error {
	return m.ReplaceFn(ctx, id, newData)
}

func (m *ImageRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

func (m *ImageRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *ImageRepositoryMock) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error) {
	return m.GetByProductIDFn(ctx, productID)
}