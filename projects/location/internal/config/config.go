package config

import (
	"github.com/spf13/viper"
	"gitlab/ArtemFed/mts-final-taxi/pkg/app"
	configLib "gitlab/ArtemFed/mts-final-taxi/pkg/config"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/projects/location/internal/repository"
	"log"
)

type Config struct {
	App              *app.Config               `mapstructure:"app"`
	Http             *http_server.Config       `mapstructure:"http"`
	Logger           *mylogger.Config          `mapstructure:"logger"`
	Postgres         *repository.Config        `mapstructure:"postgres"`
	Metrics          *metrics.Config           `mapstructure:"metrics"`
	GracefulShutdown *graceful_shutdown.Config `mapstructure:"graceful_shutdown"`
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
