package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(handler http.Handler, port string, logger *zap.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  1 * time.Minute,
		},
		logger: logger,
	}
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("HTTP server listening", zap.String("addr", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server")

	shutdownCtx, cancle := context.WithTimeout(ctx, 10*time.Second)
	defer cancle()

	return s.httpServer.Shutdown(shutdownCtx)
}
