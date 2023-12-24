package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/app"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/location/internal/config"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	ctx := context.Background()

	configPath := os.Getenv("CONFIG_LOCATION")
	if configPath == "" {
		configPath = "config/config.local.yml"
	}
	log.Println("Location config path: ", configPath)

	// Собираем конфиг приложения
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("Fail to parse Location config: ", err)
	}

	// Создаем наше приложение
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create %s app: %s", cfg.App.Name, err))
	}

	// Запускаем приложение
	// По конфигам приложение само поймет, что нужно поднять
	application.Start(ctx)
}
