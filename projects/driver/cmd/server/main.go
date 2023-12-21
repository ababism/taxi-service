package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/app"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx := context.Background()

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "config/config.local.yml"
	}
	log.Println("Driver config path: ", configPath)
	// Собираем конфиг приложения
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("Fail to parse driver config: ", err)
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
