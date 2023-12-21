package app

import (
	"context"
	"fmt"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/postgres"

	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	//"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/postgres"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"

	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	address string
	logger  *zap.Logger

	service adapters.DriverService
	//repo    adapters.UserRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	// TODO dsn
	sentry_dsn := ""
	logger, err := mylogger.InitLogger(false, sentry_dsn, "production")
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%v", cfg))

	// Чистим кэш логгера при shutdown
	graceful_shutdown.AddCallback(
		&graceful_shutdown.Callback{
			Name: "ZapLogger",
			FnCtx: func(ctx context.Context) error {
				return logger.Sync()
			},
		})

	metrics.InitOnce(cfg.Metrics, logger, metrics.AppInfo{
		Name:        cfg.App.Name,
		Environment: cfg.App.Environment,
		Version:     cfg.App.Version,
	})
	logger.Info("Init Metrics – success")

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("error while connecting to PostgreSQL DB:", zap.Error(err))
	}
	userRepo := postgres.NewDriverRepository(db)

	// TODO add shutdown Callback
	driverRepo := mongo.NewDriverRepository(logger)
	err = driverRepo.Connect(context.TODO(), cfg.Mongo, cfg.MigrationsMongo)
	if err != nil {
		logger.Fatal("error while connecting to Mongo DB:", zap.Error(err))
	}

	driverService := service.NewDriverService(userRepo)

	logger.Info("Init Driver – success")

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:     cfg,
		logger:  logger,
		service: driverService,
		//repo:    userRepo,
		address: address,
	}, nil
}
