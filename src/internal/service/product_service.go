package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/mapper"
	"github.com/4vertak/hestia/internal/repository"
	"github.com/4vertak/hestia/internal/repository/postgres"
)

type ProductService struct {
	productRepo repository.ProductRepository
	imageRepo   repository.ImageRepository
}

func NewProductService(pr repository.ProductRepository, ir repository.ImageRepository) *ProductService {
	return &ProductService{productRepo: pr, imageRepo: ir}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (s *ProductService) Create(ctx context.Context, req dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	id := uuid.New()
	product := mapper.ProductCreateRequestToModel(req, id)
	product.LastUpdateDate = timePtr(time.Now())
	if err := s.productRepo.Create(ctx, &product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	resp := mapper.ProductModelToResponse(&product)
	return &resp, nil
}

func (s *ProductService) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	err := s.productRepo.DecreaseStock(ctx, id, quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if errors.Is(err, postgres.ErrInsufficientStock) {
			return ErrInsufficientStock
		}
		return err
	}
	return nil
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrNotFound
	}
	resp := mapper.ProductModelToResponse(p)
	return &resp, nil
}

func (s *ProductService) GetAllAvailable(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.FindAllAvailable(ctx)
	if err != nil {
		return nil, err
	}
	var result []dto.ProductResponse
	for _, p := range products {
		result = append(result, mapper.ProductModelToResponse(&p))
	}
	return result, nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.productRepo.Delete(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}