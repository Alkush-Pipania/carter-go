package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/internal/modules/authentication"
	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/source"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/redis"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
)

type JWTVerifier interface {
	Verify(token string) (string, error)
}

type Container struct {
	DB                *db.Queries
	userHandler       *user.UserHandler
	collectionHandler *collection.Handler
	sourceHandler     *source.Handler
	Redis             *redis.Client
	authService       *authentication.Service
	authHandler       *authentication.Handler
}

func NewContainer(ctx context.Context, db *db.Queries, producer *rabbitmq.Producer, presigner *s3.Presigner, redis *redis.Client) *Container {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandler(userService)

	collectionRepo := collection.NewRepository(db)
	collectionService := collection.NewService(collectionRepo)
	collectionHandler := collection.NewHandler(collectionService)

	// Source module with RabbitMQ producer for embedding queue
	sourceRepo := source.NewRepository(db)
	sourceService := source.NewService(sourceRepo, producer, presigner)
	sourceHandler := source.NewHandler(sourceService)

	// Auth module
	authRepo := authentication.NewRepository(db)
	authService := authentication.NewService(authRepo, redis)
	authHandler := authentication.NewHandler(authService)

	return &Container{
		DB:                db,
		userHandler:       userHandler,
		collectionHandler: collectionHandler,
		sourceHandler:     sourceHandler,
		Redis:             redis,
		authService:       authService,
		authHandler:       authHandler,
	}
}
