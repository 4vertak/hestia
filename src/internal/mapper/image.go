package mapper

import (
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/dto"
)

func ImageModelToResponse(img *models.Image) dto.ImageResponse {
	return dto.ImageResponse{
		ID: img.ID,
	}
}