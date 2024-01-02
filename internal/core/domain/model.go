package domain

import "github.com/confluentinc/confluent-kafka-go/kafka"

// Log is a struct that represents the log model
type Log struct {
	ID               int    `json:"id" gorm:"id"`
	Level            string `json:"level" gorm:"level"`
	Message          string `json:"message" gorm:"message"`
	ResourceID       string `json:"resourceId" gorm:"resource_id"`
	Timestamp        string `json:"timestamp" gorm:"timestamp"`
	TraceID          string `json:"traceId" gorm:"trace_id"`
	SpanID           string `json:"spanId" gorm:"span_id"`
	Commit           string `json:"commit" gorm:"commit"`
	ParentResourceID string `json:"parentResourceId" gorm:"parent_resource_id"`
}

func (Log) TableName() string {
	return "log_table"
}

type LogProducer struct {
	Producer     *kafka.Producer
	Topic        string
	DeliveryChan chan kafka.Event
}
