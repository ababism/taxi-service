package config

import (
	"github.com/spf13/viper"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/app"
	configLib "gitlab.com/ArtemFed/mts-final-taxi/pkg/config"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/mytracer"
	kafkaConsumer "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/kafkaConsumer"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/scraper"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/driver_api"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/kafkaproducer"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/locationclient"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"log"
)

type Config struct {
	App              *app.Config                  `mapstructure:"app"`
	Http             *http_server.Config          `mapstructure:"http"`
	LocationClient   *locationclient.ClientConfig `mapstructure:"location_client"`
	Logger           *mylogger.Config             `mapstructure:"logger"`
	Mongo            *mongo.Config                `mapstructure:"mongo"`
	MigrationsMongo  *mongo.ConfigMigrations      `mapstructure:"migrations_mongo"`
	Metrics          *metrics.Config              `mapstructure:"metrics"`
	GracefulShutdown *graceful_shutdown.Config    `mapstructure:"graceful_shutdown"`
	KafkaReader      *kafkaConsumer.Config        `mapstructure:"kafka_reader"`
	KafkaWriter      *kafkaproducer.Config        `mapstructure:"kafka_writer"`
	Tracer           *mytracer.Config             `mapstructure:"tracer"`
	Scraper          *scraper.Config              `mapstructure:"scraper"`
	LongPoll         *driver_api.Config           `mapstructure:"long_poll"`
}

func NewConfig(filePath string, appName string) (*Config, error) {
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
	configLib.ReplaceWithEnv(&config, appName)

	return &config, nil
}
