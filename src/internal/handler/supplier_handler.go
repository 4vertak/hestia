package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/4vertak/hestia/internal/dto"
	"github.com/4vertak/hestia/internal/service"
)

type SupplierHandler struct {
	svc *service.SupplierService
}

func NewSupplierHandler(svc *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

// Create godoc
// @Summary Создать поставщика
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param body body dto.SupplierCreateRequest true "Supplier data"
// @Success 201 {object} dto.SupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/suppliers [post]
func (h *SupplierHandler) Create(c *gin.Context) {
	var req dto.SupplierCreateRequest
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

// UpdateAddress godoc
// @Summary Обновить адрес поставщика
// @Tags Suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param body body dto.SupplierUpdateAddressRequest true "New address data"
// @Success 200 {object} dto.SupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/suppliers/{id}/address [patch]
func (h *SupplierHandler) UpdateAddress(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}
	var req dto.SupplierUpdateAddressRequest
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

// Delete godoc
// @Summary Удалить поставщика
// @Tags Suppliers
// @Param id path string true "Supplier ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/suppliers/{id} [delete]
func (h *SupplierHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetAll godoc
// @Summary Получить всех поставщиков
// @Tags Suppliers
// @Produce json
// @Success 200 {array} dto.SupplierResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/suppliers [get]
func (h *SupplierHandler) GetAll(c *gin.Context) {
	res, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary Получить поставщика по ID
// @Tags Suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} dto.SupplierResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/suppliers/{id} [get]
func (h *SupplierHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id format"})
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}