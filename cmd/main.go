package main

import (
	"log"
	"log-ingestor/internal/adapters/handler"
	"log-ingestor/internal/adapters/repository"
	"log-ingestor/internal/config"
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/services"
	"log-ingestor/internal/processor"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	config.InitializeConfig()
	config.GenerateDatabaseURL()
}

var (
	healthCheckService *services.HealthCheckService
	logIngestorService *services.IngestorService
	logProducer        *domain.LogProducer
	logBulk            []domain.Log
)

func main() {

	// Connect to the database
	db, err := gorm.Open(postgres.Open(config.DatabaseConfig.DatabaseURL), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	postgresDB := repository.NewDB(db)

	// health check service
	healthCheckService = services.NewHealthCheckService(postgresDB)

	// log ingestor service
	logIngestorService = services.NewIngestorService(postgresDB)

	topic := "topic69"

	// kafka producer
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "Producer69",
		"acks":              "all",
	})

	if err != nil {
		log.Fatal(err)
	}

	// create a log producer
	logProducer = repository.NewLogProducer(kafkaProducer, topic)

	InitRoutes()
}

func InitRoutes() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// create a group public
	public := e.Group("/public")

	healthCheckHandler := handler.NewHealthCheckHandler(*healthCheckService)
	public.GET("/health", healthCheckHandler.HealthCheck)

	ingestorHandler := handler.NewIngestorHandler(*logIngestorService, logProducer)
	public.POST("/ingest", ingestorHandler.IngestLog)
	public.POST("/ingest-prepared-stmt", ingestorHandler.IngestLogWithPreparedStmt)
	public.POST("/ingest-kafka", ingestorHandler.IngestLogWithKafka)

	go processor.ConsumeLog(&logBulk, ingestorHandler)

	log.Fatal(e.Start(":1323"))
}
