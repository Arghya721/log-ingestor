package main

import (
	"log-ingestor/internal/adapters/handler"
	"log-ingestor/internal/adapters/repository"
	"log-ingestor/internal/config"
	"log-ingestor/internal/core/domain"
	"log-ingestor/internal/core/services"
	"log-ingestor/internal/processor"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	config.InitializeConfig()
	config.GenerateDatabaseURL()
	config.GetKafkaConfig()
}

// Define the Prometheus histogram
var (
	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})
)

// Custom middleware to monitor request duration
func ResponseTimeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Process request
		err := next(c)

		// Calculate duration
		duration := time.Since(startTime)

		// Record the duration in the histogram
		httpDuration.WithLabelValues(c.Request().Method, c.Path(), strconv.Itoa(c.Response().Status)).Observe(duration.Seconds())

		return err
	}
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

	// Register the Prometheus metrics
	prometheus.MustRegister(httpDuration)
	// Attach the Prometheus middleware
	e.Use(ResponseTimeMiddleware)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

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
