package http_server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Server struct {
	srv     *http.Server
	metrics *metrics
}

func New(cfg *Config) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		metrics: newMetrics(),
	}
}

func (s *Server) RegisterRoutes(r *Router) {
	s.srv.Handler = s.metricsMiddleware(r.router)
}

func (s *Server) Start() error {
	if s.srv.Handler == nil {
		return errors.Errorf("no routes have registered")
	}
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) metricsMiddleware(router *gin.Engine) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path

		lrw := newLoggingResponseWriter(w)
		router.ServeHTTP(lrw, r)
		s.metrics.observe(r.Method, strconv.Itoa(lrw.statusCode), route)
	})
}
