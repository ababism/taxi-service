package mongo

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo/migrate"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DriveRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection

	logger *zap.Logger
}

func NewDriverRepository(logger *zap.Logger) *DriveRepository {
	return &DriveRepository{logger: logger}
}

func (r *DriveRepository) Connect(ctx context.Context, cfg *MongoCfg, migrateCfg *MigrationsCfg) error {
	r.logger.Info("Connecting to mongo...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		r.logger.Error("new mongo client create error:", zap.Error(err))
		return fmt.Errorf("new mongo client create error: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		r.logger.Error("new mongo primary node connect error:", zap.Error(err))
		return fmt.Errorf("new mongo primary node connect error: %w", err)
	}

	r.client = client
	database := client.Database(cfg.Database)

	if migrateCfg.Enabled {
		migrationSvc := migrate.NewMigrationsService(r.logger, database)
		err = migrationSvc.RunMigrations(migrateCfg.Path)
		if err != nil {
			r.logger.Fatal("run migrations failed", zap.Error(err))
			return fmt.Errorf("run migrations failed")
		}
	}

	return nil
}
