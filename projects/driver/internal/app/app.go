package app

import (
	"context"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mytracer"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	kafkaConsumer "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/kafkaConsumer"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/scraper"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	driverGenerated "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/kafka_producer"
	locationClient "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client"
	locationClientGen "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client/generated"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	"go.opentelemetry.io/otel/sdk/trace"
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

	// Importing constants from driver openApi generated
	driverGenerated.ScrapeStatusesConstants()

	// Инициализируем словарь ивентов
	domain.InitTripMap()

	tp, err := mytracer.InitTracer(cfg.Tracer, cfg.App)
	if err != nil {
		return nil, err
	}
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "Opentelemetry",
			FnCtx: func(ctx context.Context) error {
				if err := tp.Shutdown(context.Background()); err != nil {
					logger.Error("Error shutting down tracer provider: %v", zap.Error(err))
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

	// Location microservice client repository layer
	client, err := locationClientGen.NewClientWithResponses(cfg.LocationClient.Uri)
	if err != nil {
		logger.Fatal("cannot initialize generated location client:", zap.Error(err))
		return nil, err
	}
	newLocationClient := locationClient.NewClient(client)
	logger.Info("Init Location client – success")

	// Kafka repository layer
	kafkaClient := kafka_producer.NewKafkaProducer(cfg.KafkaWriter)

	// Service layer
	driverService := service.NewDriverService(driverRepo, newLocationClient, kafkaClient)

	// Scraper for event calling
	scr := scraper.NewScraper(*logger, driverService)
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name:  "Data scraper stop",
			FnCtx: scr.StopFunc(),
		})

	scrapeInterval, err := cfg.Scraper.GetScrapeInterval()
	if err != nil {
		logger.Fatal("can't parse time from scraper LongPollTimeout config string:", zap.Error(err))
	}

	scr.Start(scrapeInterval)
	logger.Info("Init Scraper – success")

	// Kafka transport layer
	kafkaCumsumer := kafkaConsumer.NewKafkaConsumer(cfg.KafkaReader, driverService)
	kafkaConsumerClose := kafkaCumsumer.Start(ctx)
	graceful_shutdown.AddCallback(&graceful_shutdown.Callback{
		Name: "kafkaConsumer Close",
		FnCtx: func(ctx context.Context) error {
			return kafkaConsumerClose(kafkaCumsumer)
		},
	})
	logger.Info("Init Kafka – success")

	logger.Info("Init Driver – success")

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:            cfg,
		logger:         logger,
		service:        driverService,
		address:        address,
		tracerProvider: tp,
	}, nil
}
