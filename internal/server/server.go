package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
			WriteTimeout:  10 * time.Second,
			IdleTimeout:  1 * time.Minute,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				s.logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		s.logger.Info("shutting down server...")
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Fatal("server shutdown failed", zap.Error(err))
		}
		serverStopCtx()
	}()

	// Run the server
	s.logger.Info("server starting", zap.String("addr", s.httpServer.Addr))
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	s.logger.Info("server exited")
	return nil
}
