package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/service"
	"github.com/google/uuid"
)

type ClientHandler struct {
	svc *service.ClientService
}

func NewClientHandler(svc *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// Create godoc
// @Summary Создать нового клиента
// @Tags Clients
// @Accept json
// @Produce json
// @Param body body dto.ClientCreateRequest true "Client data"
// @Success 201 {object} dto.ClientResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/clients [post]
func (h *ClientHandler) Create(c *gin.Context) {
	var req dto.ClientCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Delete godoc
// @Summary Удалить клиента по ID
// @Tags Clients
// @Param id path string true "Client ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/clients/{id} [delete]
func (h *ClientHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}
	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// SearchByName godoc
// @Summary Поиск по имени и фамилии клиента
// @Tags Clients
// @Produce json
// @Param name query string true "Client name"
// @Param surname query string true "Client surname"
// @Success 200 {array} dto.ClientResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/clients [get]
func (h *ClientHandler) SearchByName(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")
	// страховка
	if name == "" && surname == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "name or surname is required"})
		return
	}
	res, err := h.svc.SearchByName(c.Request.Context(), name, surname)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetAll godoc
// @Summary Получить всех клиентов
// @Tags Clients
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} dto.ClientResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/clients/all [get]
func (h *ClientHandler) GetAll(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 100
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	res, err := h.svc.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateAddress godoc
// @Summary Обновить адрес клиента
// @Tags Clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param body body dto.ClientUpdateAddressRequest true "New address data"
// @Success 200 {object} dto.ClientResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/clients/{id}/address [patch]
func (h *ClientHandler) UpdateAddress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}
	var req dto.ClientUpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.svc.UpdateAddress(c.Request.Context(), id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}