package main

import (
	"gitlab/ArtemFed/mts-final-taxi/pkg/mylogger"
)

func main() {
	sentry_dsn := ""
	logger, err := mylogger.InitLogger(false, sentry_dsn, "production")
	if err != nil {
		logger.Fatal("error")
		return
	}

	logger.Info("location service: hello world")
}
