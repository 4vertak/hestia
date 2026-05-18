package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/models"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByNameAndSurname(ctx context.Context, name, surname string) ([]models.ClientWithAddress, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.ClientWithAddress, error)
	UpdateAddress(ctx context.Context, id uuid.UUID, newAddressID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
}