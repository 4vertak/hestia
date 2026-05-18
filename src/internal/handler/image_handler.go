package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/service"
)

type ImageHandler struct {
	svc *service.ImageService
}

func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// Upload godoc
// @Summary Загрузить изображение продукта
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData string true "Product ID"
// @Param file formData file true "Image file"
// @Success 201 {object} dto.ImageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/images [post]
func (h *ImageHandler) Upload(c *gin.Context) {
	productIDStr := c.PostForm("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product_id"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "missing file"})
		return
	}
	defer file.Close()

	// можно прочитать все данные
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read file"})
		return
	}

	resp, err := h.svc.Upload(c.Request.Context(), productID, buf)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Replace godoc
// @Summary Заменить изображение по ID
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Image ID"
// @Param file formData file true "New image file"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/images/{id} [put]
func (h *ImageHandler) Replace(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}

	// замена изображения: multipart form
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "missing file"})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read file"})
		return
	}

	if err := h.svc.Replace(c.Request.Context(), id, buf); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary Удалить изображение
// @Tags Images
// @Param id path string true "Image ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/images/{id} [delete]
func (h *ImageHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetByID godoc
// @Summary Получить изображение по ID
// @Tags Images
// @Produce application/octet-stream
// @Param id path string true "Image ID"
// @Success 200 {file} binary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/images/{id} [get]
func (h *ImageHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	img, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", img.Image)
}

// GetByProductID godoc
// @Summary Получить изображение продукта по ID
// @Tags Images
// @Produce application/octet-stream
// @Param productId path string true "Product ID"
// @Success 200 {file} binary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/images/product/{productId} [get]
func (h *ImageHandler) GetByProductID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid productId"})
		return
	}
	img, err := h.svc.GetByProductID(c.Request.Context(), productID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", img.Image)
}