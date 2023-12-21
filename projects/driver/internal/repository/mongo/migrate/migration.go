package migrate

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	disconnectTimeout = 30 * time.Second
)

type MigrationsService struct {
	logger migrate.Logger
	db     *mongo.Database
}

var _ migrate.Logger = (*migrationLogger)(nil)

type migrationLogger struct {
	l *zap.Logger
}

func (ml *migrationLogger) Verbose() bool {
	return true
}

func (ml *migrationLogger) Printf(format string, v ...interface{}) {
	ml.l.Info(fmt.Sprintf(format, v))
}

func NewMigrationsService(log *zap.Logger, db *mongo.Database) *MigrationsService {
	return &MigrationsService{
		logger: &migrationLogger{l: log},
		db:     db,
	}
}

func (m *MigrationsService) RunMigrations(path string) error {
	if path == "" {
		m.logger.Printf("migration was skipped")
	}
	if m.db == nil {
		return errors.New("run migration connect is not exists")
	}

	driver, err := mongodb.WithInstance(
		m.db.Client(),
		&mongodb.Config{DatabaseName: m.db.Name()},
	)
	if err != nil {
		return fmt.Errorf("cannot instantiate mongo driver: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("migration directory does not exist: %w", err)
	}

	fsrc, err := (&file.File{}).Open(path)
	if err != nil {
		return fmt.Errorf("cannot open migration source: %w", err)
	}

	instance, err := migrate.NewWithInstance("file", fsrc, "mongo", driver)
	if err != nil {
		return fmt.Errorf("new migrate instance create error: %w", err)
	}
	instance.Log = m.logger

	if err := instance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: instance.Up(): %w", err)
	}

	return nil
}
