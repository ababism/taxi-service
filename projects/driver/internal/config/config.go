package config

import (
	"github.com/spf13/viper"
	"gitlab/ArtemFed/mts-final-taxi/pkg/app"
	configLib "gitlab/ArtemFed/mts-final-taxi/pkg/config"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"log"
)

type Config struct {
	App              *app.Config               `mapstructure:"app"`
	Http             *http_server.Config       `mapstructure:"http"`
	Postgres         *Postgres                 `mapstructure:"postgres"`
	Mongo            *mongo.Config             `mapstructure:"mongo"`
	MigrationsMongo  *mongo.ConfigMigrations   `mapstructure:"migrations_mongo"`
	Metrics          *metrics.Config           `mapstructure:"metrics"`
	GracefulShutdown *graceful_shutdown.Config `mapstructure:"graceful_shutdown"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	DBName   string `mapstructure:"db-name"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl-mode"`
}

func NewConfig(filePath string) (*Config, error) {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config file: %v", err)
	}

	// Загрузка конфигурации в структуру Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error while unmarshalling config file: %v", err)
	}

	// Замена значений из переменных окружения, если они заданы
	configLib.ReplaceWithEnv(&config, "")

	return &config, nil
}
