package handler

import (
	"github.com/gin-gonic/gin"

	"flowforge-api/modules/health/v1/service"
	"flowforge-api/response"
)

type HealthHandler struct {
	service service.HealthUseCase
}

func NewHealthHandler(service service.HealthUseCase) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) Status(c *gin.Context) {
	response.Success(c, h.service.Status(c.Request.Context()), nil)
}

func (h *HealthHandler) Ping(c *gin.Context) {
	response.Success(c, h.service.Ping(), nil)
}
