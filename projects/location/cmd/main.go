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

const MainEnvName = ".env"
const AppCapsName = "LOCATION"

func init() {
	if err := godotenv.Load(MainEnvName); err != nil {
		log.Print(fmt.Sprintf("No '%s' file found", MainEnvName))
	}
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	ctx := context.Background()

	configPath := os.Getenv("CONFIG_" + AppCapsName)
	log.Println("Location config path: ", configPath)

	// Собираем конфиг приложения
	cfg, err := config.NewConfig(configPath, AppCapsName)
	if err != nil {
		log.Fatal("Fail to parse Location config: ", err)
	}

	// Создаем наше приложение
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create %s app: %s", cfg.App.Name, err))
	}

	// Запускаем приложение
	application.Start(ctx)
}
