package mongo

import (
	"context"
	"fmt"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo/models"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/service/adapters"
	"go.uber.org/zap"

	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo/migrate"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _ adapters.DriverRepository = &DriveRepository{}

type DriveRepository struct {
	client           *mongo.Client
	db               *mongo.Database
	driverCollection *mongo.Collection

	logger *zap.Logger
}

func NewDriverRepository(logger *zap.Logger) *DriveRepository {
	return &DriveRepository{logger: logger}
}

func (r *DriveRepository) Connect(ctx context.Context, cfg *Config, migrateCfg *ConfigMigrations) (func(ctx context.Context) error, error) {
	// TODO take from context
	r.logger.Info("Connecting to mongo...")
	r.logger.Info(fmt.Sprintf("mongo params: uri: %s; database: %s", cfg.Uri, cfg.Database))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		r.logger.Error("new mongo client create error:", zap.Error(err))
		return nil, fmt.Errorf("new mongo client create error: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		r.logger.Error("new mongo primary node connect error:", zap.Error(err))
		return client.Disconnect, fmt.Errorf("new mongo primary node connect error: %w", err)
	}

	r.client = client
	database := client.Database(cfg.Database)

	if migrateCfg.Enabled {
		migrationSvc := migrate.NewMigrationsService(r.logger, database)
		err = migrationSvc.RunMigrations(migrateCfg.Path)
		if err != nil {
			r.logger.Fatal("run migrations failed", zap.Error(err))
			return client.Disconnect, fmt.Errorf("run migrations failed")
		}
	}

	r.driverCollection = database.Collection(models.DriverCollectionName)

	return client.Disconnect, nil
}
