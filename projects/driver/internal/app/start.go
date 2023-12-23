package app

import (
	"context"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	requestid "github.com/sumit-tembe/gin-requestid"
	"gitlab/ArtemFed/mts-final-taxi/pkg/graceful_shutdown"
	httpServer "gitlab/ArtemFed/mts-final-taxi/pkg/http_server"
	"gitlab/ArtemFed/mts-final-taxi/pkg/metrics"
	"gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/generated"
	myHttp "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/handler/http"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"
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

	tp := trace.NewTracerProvider()

	// TODO Add tracer shutdown
	tracerMw := generated.MiddlewareFunc(otelgin.Middleware(a.cfg.App.Name, otelgin.WithTracerProvider(tp)))
	GinZapMw := generated.MiddlewareFunc(ginzap.Ginzap(a.logger, time.RFC3339, true))
	requestIdMw := generated.MiddlewareFunc(requestid.RequestID(nil))
	middlewares := []generated.MiddlewareFunc{
		tracerMw,
		GinZapMw,
		requestIdMw,
	}

	// Добавляем роуты api
	myHttp.InitHandler(router.GetRouter(), a.logger, middlewares, a.service)

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
