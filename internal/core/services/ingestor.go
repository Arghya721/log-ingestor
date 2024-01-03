package services

import (
	"encoding/json"
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/ports"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

type IngestorService struct {
	repo ports.IngestorRepository
	logs *domain.Log
}

func NewIngestorService(repo ports.IngestorRepository) *IngestorService {
	return &IngestorService{
		repo: repo,
	}
}

func (s *IngestorService) InsertLog(logRequest domain.LogRequest) error {
	err := s.repo.InsertLog(domain.Log{
		Level:            logRequest.Level,
		Message:          logRequest.Message,
		ResourceID:       logRequest.ResourceID,
		Timestamp:        logRequest.Timestamp,
		TraceID:          logRequest.TraceID,
		SpanID:           logRequest.SpanID,
		Commit:           logRequest.Commit,
		ParentResourceID: logRequest.Metadata.ParentResourceID,
	})
	if err != nil {
		log.Error("Error when inserting log", err)
		return err
	}

	return nil
}

func (s *IngestorService) InsertLogWithPreparedStmt(logRequest domain.LogRequest) error {
	err := s.repo.InsertLogWithPreparedStmt(domain.Log{
		Level:            logRequest.Level,
		Message:          logRequest.Message,
		ResourceID:       logRequest.ResourceID,
		Timestamp:        logRequest.Timestamp,
		TraceID:          logRequest.TraceID,
		SpanID:           logRequest.SpanID,
		Commit:           logRequest.Commit,
		ParentResourceID: logRequest.Metadata.ParentResourceID,
	})
	if err != nil {
		log.Error("Error when inserting log", err)
		return err
	}

	return nil
}

func (s *IngestorService) InsertLogWithKafka(logRequest domain.LogRequest, logProducer *domain.LogProducer) error {
	// convert log to byte array
	logByte, err := json.Marshal(logRequest)
	if err != nil {
		log.Error("Error when converting log to byte array", err)
		return err
	}

	// produce log to kafka
	err = logProducer.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &logProducer.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: logByte,
	},
		logProducer.DeliveryChan,
	)
	<-logProducer.DeliveryChan
	if err != nil {
		log.Error("Error when producing log to kafka", err)
		return err
	}

	return nil
}

func (s *IngestorService) InsertBulkLog(logChunk []domain.Log) {

	err := s.repo.InsertBulkLog(logChunk)
	if err != nil {
		log.Error("Error when inserting log chunk", err)
		return
	}

	return
}
