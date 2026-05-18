package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ImageRepository interface {
	Create(ctx context.Context, img *models.Image) error
	Replace(ctx context.Context, id uuid.UUID, newData []byte) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Image, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error)
}