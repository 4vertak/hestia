package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindAllAvailable(ctx context.Context) ([]models.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUpdate(ctx context.Context, id uuid.UUID) error
	UpdateImageID(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error
}