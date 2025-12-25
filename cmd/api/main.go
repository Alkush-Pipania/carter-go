package main

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/config"
	"github.com/Alkush-Pipania/carter-go/internal/app"
	"github.com/Alkush-Pipania/carter-go/internal/server"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/logger"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn := db.Init(ctx, cfg.DbUrl)
	logger.Info("Database connected")

	q := db.New(dbConn)

	// RabbitMQ Setup
	rmqConn, err := rabbitmq.NewConnection(
		rabbitmq.Config{
			URL: cfg.RabbitMQUrl,
		},
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
	})
	if err != nil {
		logger.Fatal("Failed to create RabbitMQ producer", zap.Error(err))
	}
	defer producer.Close()
	logger.Info("RabbitMQ producer created")

	redisClient, err := redisPkg.New(ctx, redisPkg.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Redis connected")

	// S3 Setup (DigitalOcean Spaces)
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

	container := app.NewContainer(ctx, q, producer, presigner, redisClient)

	router := app.NewRouter(container)

	srv := server.New(router, cfg.Port)

	logger.Info("Server starting", zap.String("port", cfg.Port))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
