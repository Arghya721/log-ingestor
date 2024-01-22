package handler

import (
	"log-ingestor/internal/core/services"

	"github.com/labstack/echo/v4"

)

type HealthCheckHandler struct {
	healthCheckService services.HealthCheckService
}

func NewHealthCheckHandler(HealthCheckService services.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthCheckService: HealthCheckService,
	}
}

func (h *HealthCheckHandler) HealthCheck(ctx echo.Context) error {
	err := h.healthCheckService.HealthCheck()
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
	}
	return ctx.JSON(200, map[string]string{
		"message": "OK",
	})
}
