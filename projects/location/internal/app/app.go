package app

import (
	"context"
	"fmt"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/config"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/repository/postgres"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/service/adapters"

	"go.uber.org/zap"
)

type App struct {
	cfg     *config.Config
	address string
	logger  *zap.Logger

	service adapters.LocationService
	//repo    adapters.LocationRepository
}

func NewApp(cfg *config.Config) (*App, error) {
	// TODO dsn
	sentryDsn := ""
	logger, err := mylogger.InitLogger(false, sentryDsn, "production")
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

	//metrics.InitOnce(cfg.Metrics, logger, metrics.AppInfo{
	//	Name:        cfg.App.Name,
	//	Environment: cfg.App.Environment,
	//	Version:     cfg.App.Version,
	//})
	//logger.Info(fmt.Sprintf("Init %s Metrics – success", cfg.App.Name))

	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logger.Fatal("error while connecting to PostgreSQL DB:", zap.Error(err))
	}
	locationRepo := postgres.NewLocalRepository(db)

	locationService := service.NewLocationService(locationRepo)

	logger.Info(fmt.Sprintf("Init %s – success", cfg.App.Name))

	address := fmt.Sprintf(":%d", cfg.Http.Port)

	return &App{
		cfg:     cfg,
		logger:  logger,
		service: locationService,
		//repo:    userRepo,
		address: address,
	}, nil
}
