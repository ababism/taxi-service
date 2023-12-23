package app

import (
	"context"
	"fmt"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/config"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/repository"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"

	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	address string
	logger  *zap.Logger

	service adapters.LocationService
}

func NewApp(cfg *config.Config) (*App, error) {
	logger, err := mylogger.InitLogger(cfg.Logger, cfg.App.Name)
	if err != nil {
		return nil, err
	}
	logger.Info("Init Logger – success")

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
	logger.Info(fmt.Sprintf("Init %s Metrics – success", cfg.App.Name))

	db, err := repository.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logger.Fatal("error while connecting to PostgreSQL DB:", zap.Error(err))
	}
	locationRepo := repository.NewLocalRepository(db)

	locationService := service.NewLocationService(locationRepo)

	logger.Info(fmt.Sprintf("Init %s – success", cfg.App.Name))

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:     cfg,
		logger:  logger,
		service: locationService,
		address: address,
	}, nil
}
