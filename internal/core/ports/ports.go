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
	InsertBulkLog(logs []domain.Log) error
}

type InternalRepository interface {
	GetLogs(logFilter domain.LogFilter) (logs domain.LogResponse, err error)
}

type IngestorService interface {
	InsertLog(log domain.LogRequest) error
	InsertLogWithPreparedStmt(log domain.LogRequest) error
	InsertLogWithKafka(log domain.LogRequest, logProducer *domain.LogProducer) error
}

type InternalService interface {
	GetLogs(logFilter domain.LogFilter) (logs domain.LogResponse, err error)
}
