package handler_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/4vertak/hestia/internal/handler"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository/mocks"
	"github.com/4vertak/hestia/internal/service"
)

func setupImageRouter(imageH *handler.ImageHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	images := r.Group("/api/v1/images")
	{
		images.POST("", imageH.Upload)
		images.PUT("/:id", imageH.Replace)
		images.GET("/:id", imageH.GetByID)
		images.GET("/product/:productId", imageH.GetByProductID)
		images.DELETE("/:id", imageH.Delete)
	}
	return r
}

func createMultipartFormData(productID uuid.UUID, fileContent []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("product_id", productID.String())
	part, _ := writer.CreateFormFile("file", "image.png")
	_, _ = io.Copy(part, bytes.NewReader(fileContent))
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUploadImage(t *testing.T) {
	productID := uuid.New()
	productRepo := &mocks.ProductRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			if id == productID {
				return &models.Product{ID: productID, ImageID: nil}, nil
			}
			return nil, nil
		},
		UpdateImageIDFn: func(ctx context.Context, pid uuid.UUID, imgID uuid.UUID) error {
			return nil
		},
	}
	imageRepo := &mocks.ImageRepositoryMock{
		CreateFn: func(ctx context.Context, img *models.Image) error {
			return nil
		},
	}
	svc := service.NewImageService(imageRepo, productRepo)
	imageH := handler.NewImageHandler(svc)
	router := setupImageRouter(imageH)

	body, contentType := createMultipartFormData(productID, []byte{0x89, 0x50, 0x4E, 0x47}) // PNG magic bytes
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/images", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetImageByID(t *testing.T) {
	imgID := uuid.New()
	imageRepo := &mocks.ImageRepositoryMock{
		GetByIDFn: func(ctx context.Context, id uuid.UUID) (*models.Image, error) {
			if id == imgID {
				return &models.Image{ID: imgID, Image: []byte{1, 2, 3}}, nil
			}
			return nil, nil
		},
	}
	productRepo := &mocks.ProductRepositoryMock{}
	svc := service.NewImageService(imageRepo, productRepo)
	imageH := handler.NewImageHandler(svc)
	router := setupImageRouter(imageH)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/images/"+imgID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, []byte{1, 2, 3}, w.Body.Bytes())
}