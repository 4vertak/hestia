package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/4vertak/hestia/internal/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	if errors.Is(err, service.ErrInsufficientStock) {
		c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}