package processor

import (
	"encoding/json"
	"log-ingestor/internal/adapters/handler"
	"log-ingestor/internal/core/domain"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

func ConsumeLog(logChunk *[]domain.Log, ingestorHandler *handler.IngestorHandler) {
	log.Info("Consuming log chunk")

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "Producer69",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to topic
	err = kafkaConsumer.SubscribeTopics([]string{"topic69"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set up a ticker to trigger insert function every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			log.Info("Inserting log chunk")
			if len(*logChunk) != 0 {
				go ingestorHandler.IngestBulkLog(logChunk)
			}
		}
	}()

	for {
		ev := kafkaConsumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Info("Message received")
			var logRequest domain.LogRequest
			err := json.Unmarshal(e.Value, &logRequest)
			if err != nil {
				log.Error("Error when unmarshalling log request", err)
			}
			log.Info("Log request", logRequest)
			*logChunk = append(*logChunk, domain.Log{
				Level:            logRequest.Level,
				Message:          logRequest.Message,
				ResourceID:       logRequest.ResourceID,
				Timestamp:        logRequest.Timestamp,
				TraceID:          logRequest.TraceID,
				SpanID:           logRequest.SpanID,
				Commit:           logRequest.Commit,
				ParentResourceID: logRequest.Metadata.ParentResourceID,
			})
		case *kafka.Error:
			log.Error("Error when consuming log", e)
		}
	}
}
