package app

import (
	"context"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	requestid "github.com/sumit-tembe/gin-requestid"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	httpServer "gitlab.com/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab.com/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	myHttp "gitlab.com/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"time"
)

// Start - Единая точка запуска приложения
func (a *App) Start(ctx context.Context) {

	go a.startHTTPServer(ctx)

	if err := graceful_shutdown.Wait(a.cfg.GracefulShutdown); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to gracefully shutdown %s app: %s", a.cfg.App.Name, err.Error()))
	} else {
		a.logger.Info("App gracefully stopped")
	}
}

func (a *App) startHTTPServer(ctx context.Context) {
	// Создаем общий роутинг http сервера
	router := httpServer.NewRouter()

	// Добавляем системные роуты
	router.WithHandleGET("/metrics", metrics.HandleFunc())

	tracerMw := generated.MiddlewareFunc(otelgin.Middleware(a.cfg.App.Name, otelgin.WithTracerProvider(a.tracerProvider)))
	GinZapMw := generated.MiddlewareFunc(ginzap.Ginzap(a.logger, time.RFC3339, true))
	requestIdMw := generated.MiddlewareFunc(requestid.RequestID(nil))
	middlewares := []generated.MiddlewareFunc{
		tracerMw,
		GinZapMw,
		requestIdMw,
	}

	longPollTimeout, err := a.cfg.LongPoll.GetLongPollTimeout()
	if err != nil {
		a.logger.Fatal("can't parse time from scraper LongPollTimeout config string:", zap.Error(err))
	}

	// Добавляем роуты api
	myHttp.InitHandler(router.GetRouter(), a.logger, middlewares, a.service, longPollTimeout)

	// Создаем сервер
	srv := httpServer.New(a.cfg.Http)
	srv.RegisterRoutes(&router)

	// Стартуем
	a.logger.Info(fmt.Sprintf("Starting %s HTTP server at %s:%d", a.cfg.App.Name, a.cfg.Http.Host, a.cfg.Http.Port))
	if err := srv.Start(); err != nil {
		a.logger.Error(fmt.Sprintf("Fail with %s HTTP server: %s", a.cfg.App.Name, err.Error()))
		graceful_shutdown.Now()
	}
}
