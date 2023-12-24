package app

import (
	"context"
	"fmt"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mytracer"
	kafkaConsumer "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/kafkaConsumer"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/kafka_producer"
	locationClient "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client"
	locationClientGen "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client/generated"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"

	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	driverGenerated "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"

	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type App struct {
	cfg            *config.Config
	address        string
	logger         *zap.Logger
	tracerProvider *trace.TracerProvider

	service adapters.DriverService
}

func NewApp(cfg *config.Config) (*App, error) {
	logger, err := mylogger.InitLogger(cfg.Logger, cfg.App.Name)
	if err != nil {
		return nil, err
	}
	logger.Info("Init Logger – success")

	startCtx := context.Background()
	ctx := zapctx.WithLogger(startCtx, logger)

	// Чистим кэш логгера при shutdown
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "ZapLogger",
			FnCtx: func(ctx context.Context) error {
				return logger.Sync()
			},
		})

	tp, err := mytracer.InitTracer(cfg.Tracer, cfg.App)
	if err != nil {
		return nil, err
	}
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "Opentelemetry",
			FnCtx: func(ctx context.Context) error {
				if err := tp.Shutdown(context.Background()); err != nil {
					log.Printf("Error shutting down tracer provider: %v", err)
					return err
				}
				return nil
			},
		})
	logger.Info("Init Tracer – success")

	// Prometheus metrics
	metrics.InitOnce(cfg.Metrics, logger, metrics.AppInfo{
		Name:        cfg.App.Name,
		Environment: cfg.App.Environment,
		Version:     cfg.App.Version,
	})
	logger.Info("Init Metrics – success")

	// MongoDB repository
	driverRepo := mongo.NewDriverRepository(logger)
	mongoDisconnect, err := driverRepo.Connect(ctx, cfg.Mongo, cfg.MigrationsMongo)
	if err != nil {
		logger.Fatal("error while connecting to Mongo DB:", zap.Error(err))
	}
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "MongoClientDisconnect",
			FnCtx: func(ctx context.Context) error {
				return mongoDisconnect(ctx)
			},
		})
	logger.Info("Mongo connect – success")

	// Location microservice client
	client, err := locationClientGen.NewClientWithResponses(cfg.LocationClient.Uri)
	if err != nil {
		log.Fatal("cannot initialize generated location client:", zap.Error(err))
		return nil, err
	}
	newLocationClient := locationClient.NewClient(client)

	// Importing constants from driver openApi generated
	driverGenerated.ScrapeStatusesConstants()
	kafkaClient := kafka_producer.NewKafkaProducer(cfg.KafkaWriter)

	driverService := service.NewDriverService(driverRepo, newLocationClient, kafkaClient)

	logger.Info("Init Driver – success")

	kafkaCumsumer := kafkaConsumer.NewKafkaConsumer(cfg.KafkaReader, driverService)
	kafkaConsumerClose := kafkaCumsumer.Start(ctx)
	graceful_shutdown.AddCallback(&graceful_shutdown.Callback{
		Name: "MongoClientDisconnect",
		FnCtx: func(ctx context.Context) error {
			return kafkaConsumerClose(kafkaCumsumer)
		},
	})

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:            cfg,
		logger:         logger,
		service:        driverService,
		address:        address,
		tracerProvider: tp,
	}, nil
}
