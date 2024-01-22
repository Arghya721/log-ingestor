package handler

import (
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/services"

	"github.com/labstack/echo/v4"
)

type IngestorHandler struct {
	ingestorService services.IngestorService
	logProducer     *domain.LogProducer
}

func NewIngestorHandler(IngestorService services.IngestorService, logProducer *domain.LogProducer) *IngestorHandler {
	return &IngestorHandler{
		ingestorService: IngestorService,
		logProducer:     logProducer,
	}
}

func (h *IngestorHandler) IngestLog(ctx echo.Context) error {
	var log domain.LogRequest
	err := ctx.Bind(&log)
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"message": "Bad Request",
		})
	}
	err = h.ingestorService.InsertLog(log)
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
	}
	return ctx.JSON(200, map[string]string{
		"message": "OK",
	})
}

func (h *IngestorHandler) IngestLogWithPreparedStmt(ctx echo.Context) error {
	var log domain.LogRequest
	err := ctx.Bind(&log)
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"message": "Bad Request",
		})
	}
	err = h.ingestorService.InsertLogWithPreparedStmt(log)
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
	}
	return ctx.JSON(200, map[string]string{
		"message": "OK",
	})
}

func (h *IngestorHandler) IngestLogWithKafka(ctx echo.Context) error {
	var log domain.LogRequest
	err := ctx.Bind(&log)
	if err != nil {
		return ctx.JSON(400, map[string]string{
			"message": "Bad Request",
		})
	}
	err = h.ingestorService.InsertLogWithKafka(log, h.logProducer)
	if err != nil {
		return ctx.JSON(500, map[string]string{
			"message": "Internal Server Error",
		})
	}
	return ctx.JSON(200, map[string]string{
		"message": "OK",
	})
}

func (h *IngestorHandler) IngestBulkLog(logRequest []domain.Log) {
	h.ingestorService.InsertBulkLog(logRequest)

	return
}
