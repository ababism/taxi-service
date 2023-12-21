package config

import (
	"github.com/spf13/viper"
	"gitlab/ArtemFed/mts-final-taxi/pkg/app"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	"gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/repository/mongo"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	replaceWithEnv(&config, "")

	return &config, nil
}

// replaceWithEnv заменяет все значения структуры на соответствующие значения из переменных окружения
func replaceWithEnv(config interface{}, prefix string) {
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		fieldName := v.Type().Field(i).Name

		if prefix != "" {
			fieldName = prefix + "_" + fieldName
		}

		if fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct {
			// Рекурсивный вызов для вложенных структур, если не nil
			if !field.IsNil() {
				replaceWithEnv(field.Interface(), fieldName)
			}
		} else {
			envName := strings.ToUpper(fieldName)
			envValue := os.Getenv(envName)

			// Замена значения, если переменная окружения задана
			if envValue != "" {
				log.Printf(
					"The configuration value of %s has been overwritten from %s to %s\n",
					envName,
					field.String(),
					envValue,
				)
				setField(&field, envValue)
			}
		}
	}
}

// setField устанавливает значение поля с учетом его типа
func setField(field *reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			field.SetFloat(floatValue)
		}
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			field.SetBool(boolValue)
		}
	default:
		return
	}
}
