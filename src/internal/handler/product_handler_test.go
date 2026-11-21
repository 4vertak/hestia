package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/handler"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository/mocks"
	"github.com/4vertak/hestia/internal/repository/postgres"
	"github.com/4vertak/hestia/internal/service"
)

func setupProductRouter(productH *handler.ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	products := r.Group("/api/v1/products")
	{
		products.POST("", productH.Create)
		products.PATCH("/:id/decrease", productH.DecreaseStock)
		products.GET("/:id", productH.GetByID)
		products.GET("", productH.GetAllAvailable)
		products.DELETE("/:id", productH.Delete)
	}
	return r
}

func TestCreateProduct(t *testing.T) {
	productRepo := &mocks.ProductRepositoryMock{
		CreateFn: func(ctx context.Context, p *models.Product) error {
			return nil
		},
	}
	imageRepo := &mocks.ImageRepositoryMock{}
	svc := service.NewProductService(productRepo, imageRepo)
	productH := handler.NewProductHandler(svc)
	router := setupProductRouter(productH)

	reqBody := dto.ProductCreateRequest{
		Name:           "Microwave",
		Category:       "Kitchen",
		Price:          299.99,
		AvailableStock: 10,
		SupplierID:     uuid.New(),
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDecreaseStockInsufficient(t *testing.T) {
	productID := uuid.New()
	productRepo := &mocks.ProductRepositoryMock{
		DecreaseStockFn: func(ctx context.Context, id uuid.UUID, quantity int) error {
			return postgres.ErrInsufficientStock
		},
	}
	imageRepo := &mocks.ImageRepositoryMock{}
	svc := service.NewProductService(productRepo, imageRepo)
	productH := handler.NewProductHandler(svc)
	router := setupProductRouter(productH)

	reqBody := dto.ProductDecreaseRequest{Quantity: 999}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/products/"+productID.String()+"/decrease", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}