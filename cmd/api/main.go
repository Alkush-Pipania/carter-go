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

	log := logger.Get() // Get the initialized logger

	log.Info("Starting application",
		zap.String("env", cfg.Env),
		zap.String("log_level", cfg.LogLevel),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn := db.Init(ctx, cfg.DbUrl)
	log.Info("Database connected")

	q := db.New(dbConn)

	// RabbitMQ Setup
	rmqConn, err := rabbitmq.NewConnection(
		rabbitmq.Config{
			URL: cfg.RabbitMQUrl,
		},
	)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rmqConn.Close()
	log.Info("RabbitMQ connected")

	producer, err := rabbitmq.NewProducer(rmqConn, rabbitmq.ProducerConfig{
		Exchange:     "carter.embedding",
		ExchangeType: "direct",
		Durable:      true,
	})
	if err != nil {
		log.Fatal("Failed to create RabbitMQ producer", zap.Error(err))
	}
	defer producer.Close()
	log.Info("RabbitMQ producer created")

	redisClient, err := redisPkg.New(ctx, redisPkg.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	log.Info("Redis connected")

	// S3 Setup (DigitalOcean Spaces)
	s3Client, err := s3.NewClient(ctx, s3.ClientConfig{
		Region:     cfg.DORegion,
		Endpoint:   cfg.DOEndpoint,
		AccessKey:  cfg.DOAccessKey,
		SecretKey:  cfg.DOSecretKey,
		BucketName: cfg.DOBucket,
	})
	if err != nil {
		log.Fatal("Failed to create S3 client", zap.Error(err))
	}
	log.Info("S3 client initialized (DigitalOcean Spaces)",
		zap.String("region", cfg.DORegion),
		zap.String("bucket", cfg.DOBucket))
	presigner := s3.NewPresigner(s3Client, s3.PresignerConfig{
		ExpiryMinutes: cfg.PresignExpiry,
	}, log)
	log.Info("S3 presigner created",
		zap.Int("expiry_minutes", cfg.PresignExpiry))

	// Wiring dependencies using NewApp which returns the router (http.Handler)
	// We use the simpler wiring function we created.
	// Note: NewContainer was renamed/refactored to NewApp in internal/app/container.go
	router := app.NewApp(ctx, q, producer, presigner, redisClient)

	// Initialize Server with the new signature that accepts logger
	srv := server.New(router, cfg.Port, log)

	log.Info("Server initialized", zap.String("port", cfg.Port))

	// Run the server (blocking call that handles graceful shutdown)
	if err := srv.Run(); err != nil {
		log.Fatal("Server exited with error", zap.Error(err))
	}
}
