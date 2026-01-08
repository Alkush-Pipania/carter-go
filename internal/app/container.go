package app

import (
	"context"
	"net/http"

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

// NewApp initializes dependencies and returns the router (http.Handler)
func NewApp(ctx context.Context, db *db.Queries, producer *rabbitmq.Producer, presigner *s3.Presigner, redis *redis.Client) http.Handler {
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

	// We can reuse NewRouter's logic here or call a modified NewRouter
	// Since NewRouter depended on Container, we will inline the routing logic or refactor NewRouter.
	// For better separation, let's keep NewRouter but make it accept the handlers.
	// But given the task is to refactor container, let's make this function return the handler directly.

	return NewRouter(userHandler, collectionHandler, sourceHandler, uploadHandler, authHandler, authService)
}
