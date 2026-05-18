package service

import (
	"context"
	"database/sql"
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/mapper"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository"
)

type ImageService struct {
	imageRepo   repository.ImageRepository
	productRepo repository.ProductRepository
}

func NewImageService(ir repository.ImageRepository, pr repository.ProductRepository) *ImageService {
	return &ImageService{imageRepo: ir, productRepo: pr}
}

func (s *ImageService) Upload(ctx context.Context, productID uuid.UUID, reader io.Reader) (*dto.ImageResponse, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil || product == nil {
		return nil, ErrNotFound
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// если у продукта уже было изображение, удаляем
	if product.ImageID != nil {
		_ = s.imageRepo.Delete(ctx, *product.ImageID)
	}

	imgID := uuid.New()
	img := models.Image{
		ID:    imgID,
		Image: data,
	}
	if err := s.imageRepo.Create(ctx, &img); err != nil {
		return nil, err
	}
	if err := s.productRepo.UpdateImageID(ctx, productID, imgID); err != nil {
		return nil, err
	}

	resp := mapper.ImageModelToResponse(&img)
	resp.ProductID = productID
	return &resp, nil
}

func (s *ImageService) Replace(ctx context.Context, imageID uuid.UUID, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = s.imageRepo.Replace(ctx, imageID, data)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *ImageService) Delete(ctx context.Context, imageID uuid.UUID) error {
	err := s.imageRepo.Delete(ctx, imageID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *ImageService) GetByID(ctx context.Context, imageID uuid.UUID) (*models.Image, error) {
	img, err := s.imageRepo.GetByID(ctx, imageID)
	if err != nil || img == nil {
		return nil, ErrNotFound
	}
	return img, nil
}

func (s *ImageService) GetByProductID(ctx context.Context, productID uuid.UUID) (*models.Image, error) {
	img, err := s.imageRepo.GetByProductID(ctx, productID)
	if err != nil || img == nil {
		return nil, ErrNotFound
	}
	return img, nil
}