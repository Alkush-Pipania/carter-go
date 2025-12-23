package main

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/config"
	"github.com/Alkush-Pipania/carter-go/internal/app"
	"github.com/Alkush-Pipania/carter-go/internal/server"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/logger"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
	"github.com/Alkush-Pipania/carter-go/pkg/utils"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn := db.Init(ctx, cfg.DbUrl)
	logger.Info("Database connected")

	q := db.New(dbConn)

	jwt := utils.NewJwtservice(cfg.JwtSecret)

	// RabbitMQ Setup
	rmqConn, err := rabbitmq.NewConnection(
		rabbitmq.DefaultConfig(cfg.RabbitMQUrl),
		logger.Get(),
	)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rmqConn.Close()
	logger.Info("RabbitMQ connected")

	producer, err := rabbitmq.NewProducer(rmqConn, rabbitmq.ProducerConfig{
		Exchange:     "carter.embedding",
		ExchangeType: "direct",
		Durable:      true,
	}, logger.Get())
	if err != nil {
		logger.Fatal("Failed to create RabbitMQ producer", zap.Error(err))
	}
	defer producer.Close()
	logger.Info("RabbitMQ producer created")

	// S3 Setup
	s3Client, err := s3.NewClient(ctx, s3.ClientConfig{
		Region:     cfg.AWSRegion,
		BucketName: cfg.S3BucketName,
	}, logger.Get())
	if err != nil {
		logger.Fatal("Failed to create S3 client", zap.Error(err))
	}
	logger.Info("S3 client initialized",
		zap.String("region", cfg.AWSRegion),
		zap.String("bucket", cfg.S3BucketName))

	presigner := s3.NewPresigner(s3Client, s3.PresignerConfig{
		ExpiryMinutes: cfg.PresignExpiry,
	}, logger.Get())
	logger.Info("S3 presigner created",
		zap.Int("expiry_minutes", cfg.PresignExpiry))

	container := app.NewContainer(ctx, q, jwt, producer, presigner)

	router := app.NewRouter(container)

	srv := server.New(router, cfg.Port)

	logger.Info("Server starting", zap.String("port", cfg.Port))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
