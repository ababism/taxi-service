package config

import (
	"os"

	"gitlab/ArtemFed/mts-final-taxi/pkg/app"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App              *app.Info                 `yaml:"app"`
	Http             *http_server.Config       `yaml:"http"`
	Postgres         *Postgres                 `yaml:"postgres"`
	Mongo            *mongo.MongoCfg           `yaml:"mongo"`
	MigrationsMongo  *mongo.MigrationsCfg      `yaml:"migrations_mongo"`
	Metrics          *metrics.Config           `yaml:"metrics"`
	GracefulShutdown *graceful_shutdown.Config `yaml:"graceful_shutdown"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	DBName   string `yaml:"db-name"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl-mode"`
}

// TODO из флагов
func NewConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
