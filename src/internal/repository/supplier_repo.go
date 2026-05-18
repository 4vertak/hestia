package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *models.Supplier) error
	UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]models.SupplierWithAddress, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.SupplierWithAddress, error)
}