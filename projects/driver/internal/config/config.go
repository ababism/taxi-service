package config

import (
	"github.com/spf13/viper"
	"gitlab/ArtemFed/mts-final-taxi/pkg/app"
	configLib "gitlab/ArtemFed/mts-final-taxi/pkg/config"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
	"gitlab/ArtemFed/mts-final-taxi/pkg/mytracer"
	kafkaConsumer "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/kafkaConsumer"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/daemons/scraper"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http/driver_api"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/kafka_producer"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/location_client"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"log"
)

type Config struct {
	App              *app.Config                   `mapstructure:"app"`
	Http             *http_server.Config           `mapstructure:"http"`
	LocationClient   *location_client.ClientConfig `mapstructure:"location_client"`
	Logger           *mylogger.Config              `mapstructure:"logger"`
	Mongo            *mongo.Config                 `mapstructure:"mongo"`
	MigrationsMongo  *mongo.ConfigMigrations       `mapstructure:"migrations_mongo"`
	Metrics          *metrics.Config               `mapstructure:"metrics"`
	GracefulShutdown *graceful_shutdown.Config     `mapstructure:"graceful_shutdown"`
	KafkaReader      *kafkaConsumer.Config         `mapstructure:"kafka_reader"`
	KafkaWriter      *kafka_producer.Config        `mapstructure:"kafka_writer"`
	Tracer           *mytracer.Config              `mapstructure:"tracer"`
	Scraper          *scraper.Config               `mapstructure:"scraper"`
	LongPoll         *driver_api.Config            `mapstructure:"long_poll"`
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
