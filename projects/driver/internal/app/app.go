package app

import (
	"context"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/mytracer"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	kafkaConsumer "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/kafkaConsumer"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/scraper"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/domain"
	driverGenerated "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/kafkaproducer"
	locationClient "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/locationclient"
	locationClientGen "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/locationclient/generated"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

type App struct {
	cfg            *config.Config
	address        string
	logger         *zap.Logger
	tracerProvider *trace.TracerProvider
	service        adapters.DriverService
}

func NewApp(cfg *config.Config) (*App, error) {
	startCtx := context.Background()

	// Инициализируем логгер
	logger, err := mylogger.InitLogger(cfg.Logger, cfg.App.Name)
	if err != nil {
		return nil, err
	}
	// Чистим кэш логгера при shutdown
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "ZapLogger",
			FnCtx: func(ctx context.Context) error {
				return logger.Sync()
			},
		})
	logger.Info("Init Logger – success")

	ctx := zapctx.WithLogger(startCtx, logger)

	// Importing constants from driver openApi generated
	driverGenerated.ScrapeStatusesConstants()

	// Инициализируем словарь ивентов
	domain.InitTripMap()

	// Инициализируем трейсер
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

	// Инициализируем Prometheus
	metrics.InitOnce(cfg.Metrics, logger, metrics.AppInfo{
		Name:        cfg.App.Name,
		Environment: cfg.App.Environment,
		Version:     cfg.App.Version,
	})
	logger.Info("Init Metrics – success")

	driverRepo := mongo.NewDriverRepository(logger)

	// Инициализируем MongoDB
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

	// Инициализируем Location microservice client repository layer
	client, err := locationClientGen.NewClientWithResponses(cfg.LocationClient.Uri)
	if err != nil {
		logger.Fatal("cannot initialize generated location client:", zap.Error(err))
		return nil, err
	}
	newLocationClient := locationClient.NewClient(client)
	logger.Info("Init Location client – success")

	// Kafka repository layer
	kafkaClient := kafkaproducer.NewKafkaProducer(cfg.KafkaWriter)

	// Service layer
	driverService := service.NewDriverService(driverRepo, newLocationClient, kafkaClient)

	logger.Info(fmt.Sprintf("Init %s – success", cfg.App.Name))

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

	// Kafka transport layer (тихо шифером шурша... крыша едет не спеша)
	kakafkaCumsumer := kafkaConsumer.NewKafkaConsumer(cfg.KafkaReader, driverService)
	kakafkaCumsumerClose := kakafkaCumsumer.Start(ctx)
	graceful_shutdown.AddCallback(&graceful_shutdown.Callback{
		Name: "KakafkaCumsumer Close",
		FnCtx: func(ctx context.Context) error {
			return kakafkaCumsumerClose(kakafkaCumsumer)
		},
	})
	logger.Info("Init Kafka – success")

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:            cfg,
		logger:         logger,
		service:        driverService,
		address:        address,
		tracerProvider: tp,
	}, nil
}
