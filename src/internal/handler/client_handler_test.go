package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/handler"
	"github.com/4vertak/hestia/internal/models"
	"github.com/4vertak/hestia/internal/repository/mocks"
	"github.com/4vertak/hestia/internal/service"
)

func setupRouter(clientH *handler.ClientHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	clients := r.Group("/api/v1/clients")
	{
		clients.POST("", clientH.Create)
		clients.GET("", clientH.SearchByName)
		clients.GET("/all", clientH.GetAll)
		clients.PATCH("/:id/address", clientH.UpdateAddress)
		clients.DELETE("/:id", clientH.Delete)
	}
	return r
}

func TestCreateClient(t *testing.T) {
	// arrange
	addressRepo := &mocks.AddressRepositoryMock{
		CreateFn: func(ctx context.Context, addr *models.Address) error {
			addr.ID = uuid.New()
			return nil
		},
	}

	clientRepo := &mocks.ClientRepositoryMock{
		CreateFn: func(ctx context.Context, c *models.Client) error {
			return nil
		},
	}

	svc := service.NewClientService(clientRepo, addressRepo)
	clientH := handler.NewClientHandler(svc)
	router := setupRouter(clientH)

	reqBody := dto.ClientCreateRequest{
		Name:    "John",
		Surname: "Doe",
		Gender:  strPtr("M"),
		Address: dto.AddressRequest{
			Country: "USA",
			City:    "New York",
			Street:  "5th Avenue",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.ClientResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "John", resp.Name)
	assert.Equal(t, "USA", resp.Address.Country)
}

func TestCreateClientBadRequest(t *testing.T) {
	clientH := handler.NewClientHandler(&service.ClientService{})
	router := setupRouter(clientH)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteClientNotFound(t *testing.T) {
	clientRepo := &mocks.ClientRepositoryMock{
		DeleteFn: func(ctx context.Context, id uuid.UUID) error {
			return sql.ErrNoRows
		},
	}
	svc := service.NewClientService(clientRepo, &mocks.AddressRepositoryMock{})
	clientH := handler.NewClientHandler(svc)
	router := setupRouter(clientH)

	id := uuid.New()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/clients/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Сервис оборачивает sql.ErrNoRows в ErrNotFound, хендлер вернет 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Вспомогательная функция
func strPtr(s string) *string {
	return &s
}