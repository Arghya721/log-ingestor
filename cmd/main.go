package main

import (
	"log-ingestor/internal/adapters/handler"
	"log-ingestor/internal/adapters/repository"
	"log-ingestor/internal/config"
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/services"
	"log-ingestor/internal/processor"

	"github.com/labstack/gommon/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	config.InitializeConfig()
	config.GenerateDatabaseURL()
	config.GetKafkaConfig()
}

var (
	healthCheckService *services.HealthCheckService
	logIngestorService *services.IngestorService
	internalService    *services.InternalService
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

	// internal api service
	internalService = services.NewInternalService(postgresDB)

	// kafka producer
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaConfig.KafkaConnectionURL,
		"client.id":         "xyz",
		"acks":              "all",
	})

	if err != nil {
		log.Error(err)
	}

	// create a log producer
	logProducer = repository.NewLogProducer(kafkaProducer)

	InitRoutes()
}

func InitRoutes() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// create a group public
	public := e.Group("/public")

	// create an internal group
	internal := e.Group("/internal")

	healthCheckHandler := handler.NewHealthCheckHandler(*healthCheckService)
	public.GET("/health", healthCheckHandler.HealthCheck)

	ingestorHandler := handler.NewIngestorHandler(*logIngestorService, logProducer)
	public.POST("/ingest", ingestorHandler.IngestLog)
	public.POST("/ingest-prepared-stmt", ingestorHandler.IngestLogWithPreparedStmt)
	public.POST("/ingest-kafka", ingestorHandler.IngestLogWithKafka)

	internalHandler := handler.NewInternalHandler(*internalService)
	internal.GET("/logs", internalHandler.GetLogs)

	go processor.ConsumeLog(&logBulk, ingestorHandler)

	log.Fatal(e.Start(":1323"))
}
