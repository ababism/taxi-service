package main

import (
	"context"
	"log"

	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/app"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/config"
)

func main() {
	ctx := context.Background()

	configPath := "projects/driver/config/config.local.yml"
	// Собираем конфиг приложения
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("Fail to parse config: ", err)
	}

	// Создаем наше приложение
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal("Fail to create app: ", err)
	}

	// Запускаем приложение
	// По конфигам приложение само поймет, что нужно поднять
	application.Start(ctx)
}
