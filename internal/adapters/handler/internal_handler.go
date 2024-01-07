package handler

import (
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/services"
	"strconv"

	"github.com/labstack/echo"
)

type InternalHandler struct {
	internalService services.InternalService
}

func NewInternalHandler(internalService services.InternalService) *InternalHandler {
	return &InternalHandler{
		internalService: internalService,
	}
}

func (h *InternalHandler) GetLogs(ctx echo.Context) error {
	// get log filter from query params
	var logFilter domain.LogFilter

	// get page and limit from query params
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"message": "Invalid page value",
		})
	}
	logFilter.Page = page

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"message": "Invalid limit value",
		})
	}
	logFilter.Limit = limit

	// get logs
	logs, err := h.internalService.GetLogs(logFilter)
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(200, logs)
}
