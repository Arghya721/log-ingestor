package processor

import (
	"encoding/json"
	"log-ingestor/internal/adapters/handler"
	"log-ingestor/internal/core/domain"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

func ConsumeLog(logChunk *[]domain.Log, ingestorHandler *handler.IngestorHandler) {
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

	for {
		// Read 1000 messages at a time
		for i := 0; i < 1000; i++ {
			ev := kafkaConsumer.Poll(1)
			switch e := ev.(type) {
			case *kafka.Message:
				log.Info("Message received")
				var logRequest domain.LogRequest
				err := json.Unmarshal(e.Value, &logRequest)
				if err != nil {
					log.Error("Error when unmarshalling log request", err)
				}
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

		// Insert the batch of 1000 messages into the database
		if len(*logChunk) != 0 {
			var logRequest []domain.Log
			logRequest = *logChunk

			// Reset log chunk
			*logChunk = []domain.Log{}
			log.Info("Inserting log chunk")
			ingestorHandler.IngestBulkLog(logRequest)
		}
	}
}
