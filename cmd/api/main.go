package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alkush-Pipania/carter-go/config"
	"github.com/Alkush-Pipania/carter-go/internal/app"
	"github.com/Alkush-Pipania/carter-go/internal/server"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/logger"
	redisPkg "github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadEnv()

	logger.Init(logger.Config{
		Env:   cfg.Env,
		Level: cfg.LogLevel,
	})
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("env", cfg.Env),
		zap.String("log_level", cfg.LogLevel),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbConn := db.Init(ctx, cfg.DbUrl)
	logger.Info("Database connected")

	q := db.New(dbConn)

	redisClient, err := redisPkg.New(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connected")

	s3Client, err := s3.NewClient(ctx, s3.ClientConfig{
		Region:     cfg.DORegion,
		Endpoint:   cfg.DOEndpoint,
		AccessKey:  cfg.DOAccessKey,
		SecretKey:  cfg.DOSecretKey,
		BucketName: cfg.DOBucket,
	})
	if err != nil {
		logger.Fatal("Failed to create S3 client", zap.Error(err))
	}
	logger.Info("S3 client initialized (DigitalOcean Spaces)",
		zap.String("region", cfg.DORegion),
		zap.String("bucket", cfg.DOBucket))
	presigner := s3.NewPresigner(s3Client, s3.PresignerConfig{
		ExpiryMinutes: cfg.PresignExpiry,
	}, logger.Get())
	logger.Info("S3 presigner created",
		zap.Int("expiry_minutes", cfg.PresignExpiry))

	container, err := app.NewContainer(ctx, q, presigner, redisClient, cfg)
	if err != nil {
		logger.Fatal("Failed to connect", zap.Error(err))

	}

	router := app.NewRouter(container)

	srv := server.New(router, cfg.Port, logger.Get())
	srv.Start()

	<-ctx.Done() // wait for the signal
	logger.Info("shutdown signal received")

	// 1. stop the http server
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := container.Shutdown(shutdownCtx); err != nil {
		logger.Error("dependencies shutdown failed", zap.Error(err))
	}

	// Shutdown done
	logger.Info("graceful shutdown complete")

}
