package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/config"
	"github.com/Alkush-Pipania/carter-go/internal/modules/authentication"
	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/source"
	"github.com/Alkush-Pipania/carter-go/internal/modules/upload"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
	"github.com/Alkush-Pipania/carter-go/pkg/validation"
	"github.com/rabbitmq/amqp091-go"
)

type JWTVerifier interface {
	Verify(token string) (string, error)
}

type Container struct {
	DB                *db.Queries
	userHandler       *user.UserHandler
	collectionHandler *collection.Handler
	sourceHandler     *source.Handler
	uploadHandler     *upload.Handler
	Redis             *redis.Client
	authService       *authentication.Service
	authHandler       *authentication.Handler
	RMQConn           *amqp091.Connection
}

func NewContainer(ctx context.Context, db *db.Queries, presigner *s3.Presigner, redis *redis.Client, cfg *config.Config) (*Container, error) {
	// Shared validator instance for all handlers
	validator := validation.NewValidator()

	rmqConn, err := setupRabbitMQ(&cfg.RabbitMQ)
	if err != nil {
		return nil, err
	}

	pbh, err := rabbitmq.NewPublisher(rmqConn, cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.RoutingKey)
	if err != nil {
		return nil, err
	}
	delPbh, err := rabbitmq.NewPublisher(rmqConn, cfg.RabbitMQ.ExchangeName, cfg.RabbitMQ.DeleteRoutingKey)
	if err != nil {
		return nil, err
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandler(userService)

	collectionRepo := collection.NewRepository(db)
	collectionService := collection.NewService(collectionRepo)
	collectionHandler := collection.NewHandler(collectionService, validator)

	// Source module with RabbitMQ producer for embedding queue
	sourceRepo := source.NewRepository(db)
	sourceService := source.NewService(sourceRepo, pbh, delPbh)
	sourceHandler := source.NewHandler(sourceService, validator)

	// Upload module with S3 presigner and RabbitMQ producer for file uploads
	uploadRepo := upload.NewRepository(db)
	uploadService := upload.NewService(uploadRepo, presigner, pbh)
	uploadHandler := upload.NewHandler(uploadService, validator)

	// Auth module
	authRepo := authentication.NewRepository(db)
	authService := authentication.NewService(authRepo, redis, userService)
	authHandler := authentication.NewHandler(authService, validator)

	return &Container{
		DB:                db,
		userHandler:       userHandler,
		collectionHandler: collectionHandler,
		sourceHandler:     sourceHandler,
		uploadHandler:     uploadHandler,
		Redis:             redis,
		authService:       authService,
		authHandler:       authHandler,
		RMQConn:           rmqConn,
	}, nil
}

func setupRabbitMQ(cfg *config.RabbitMQConfig) (*amqp091.Connection, error) {
	rmqpConn, err := rabbitmq.NewConn(cfg)
	if err != nil {
		return nil, err
	}
	if err := rabbitmq.SetupTopology(rmqpConn, cfg); err != nil {
		return nil, err
	}
	return rmqpConn, nil
}

func (c *Container) Shutdown(ctx context.Context) error {
	if c.RMQConn != nil {
		_ = c.RMQConn.Close()
	}
	return nil
}
