package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type AddressRepository interface {
	Create(ctx context.Context, addr *models.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Address, error)
	Update(ctx context.Context, addr *models.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}