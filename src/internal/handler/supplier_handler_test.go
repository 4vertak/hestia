package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/handler"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository/mocks"
	"github.com/4vertak/hestia/internal/service"
)

func setupSupplierRouter(supplierH *handler.SupplierHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	suppliers := r.Group("/api/v1/suppliers")
	{
		suppliers.POST("", supplierH.Create)
		suppliers.GET("", supplierH.GetAll)
		suppliers.GET("/:id", supplierH.GetByID)
		suppliers.PATCH("/:id/address", supplierH.UpdateAddress)
		suppliers.DELETE("/:id", supplierH.Delete)
	}
	return r
}

func TestCreateSupplier(t *testing.T) {
	addressRepo := &mocks.AddressRepositoryMock{
		CreateFn: func(ctx context.Context, addr *models.Address) error {
			return nil
		},
	}
	supplierRepo := &mocks.SupplierRepositoryMock{
		CreateFn: func(ctx context.Context, s *models.Supplier) error {
			return nil
		},
	}
	svc := service.NewSupplierService(supplierRepo, addressRepo)
	supplierH := handler.NewSupplierHandler(svc)
	router := setupSupplierRouter(supplierH)

	reqBody := dto.SupplierCreateRequest{
		Name:    "Best Appliances Ltd",
		Address: dto.AddressRequest{Country: "UK", City: "London", Street: "Oxford Street"},
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}