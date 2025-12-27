package app

import (
	"context"

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
}

func NewContainer(ctx context.Context, db *db.Queries, producer *rabbitmq.Producer, presigner *s3.Presigner, redis *redis.Client) *Container {
	// Shared validator instance for all handlers
	validator := validation.NewValidator()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandler(userService)

	collectionRepo := collection.NewRepository(db)
	collectionService := collection.NewService(collectionRepo)
	collectionHandler := collection.NewHandler(collectionService, validator)

	// Source module with RabbitMQ producer for embedding queue
	sourceRepo := source.NewRepository(db)
	sourceService := source.NewService(sourceRepo, producer)
	sourceHandler := source.NewHandler(sourceService, validator)

	// Upload module with S3 presigner and RabbitMQ producer for file uploads
	uploadRepo := upload.NewRepository(db)
	uploadService := upload.NewService(uploadRepo, presigner, producer)
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
	}
}
