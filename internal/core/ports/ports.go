package ports

import (
	"log-ingestor/internal/core/domain"
)

type CheckHealthRepository interface {
	HealthCheck() error
}

type IngestorRepository interface {
	InsertLog(log domain.Log) error
	InsertLogWithPreparedStmt(log domain.Log) error
	InsertBulkLog(logs *[]domain.Log) error
}

type IngestorService interface {
	InsertLog(log domain.LogRequest) error
	InsertLogWithPreparedStmt(log domain.LogRequest) error
	InsertLogWithKafka(log domain.LogRequest, logProducer *domain.LogProducer) error
}
