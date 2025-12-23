package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/source"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
	"github.com/Alkush-Pipania/carter-go/pkg/rabbitmq"
	"github.com/Alkush-Pipania/carter-go/pkg/s3"
)

type JWTVerifier interface {
	Verify(token string) (string, error)
}

type Container struct {
	DB                *db.Queries
	userHandler       *user.UserHandler
	collectionHandler *collection.Handler
	jwt               JWTVerifier
	sourceHandler     *source.Handler
}

func NewContainer(ctx context.Context, db *db.Queries, jwt JWTVerifier, producer *rabbitmq.Producer, presigner *s3.Presigner) *Container {
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

	return &Container{
		DB:                db,
		userHandler:       userHandler,
		collectionHandler: collectionHandler,
		sourceHandler:     sourceHandler,
		jwt:               jwt,
	}
}
