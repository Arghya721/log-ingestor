package repository

import (
	"context"
	"log-ingestor/internal/config"
	"log-ingestor/internal/core/domain"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

func NewLogProducer(p *kafka.Producer) *domain.LogProducer {
	return &domain.LogProducer{
		Producer:     p,
		Topic:        config.KafkaConfig.KafkaTopic,
		DeliveryChan: make(chan kafka.Event, 10000),
	}
}

func (repo *DB) HealthCheck() error {
	db, err := repo.db.WithContext(context.Background()).DB()
	if err != nil {
		return err
	}
	log.Info("Pinging the database")
	return db.Ping()
}
