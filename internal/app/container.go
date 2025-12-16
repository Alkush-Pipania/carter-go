package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/internal/modules/collection"
	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

type Container struct {
	DB                *db.Queries
	userHandler       *user.UserHandler
	collectionHandler *collection.Handler
}

func NewContainer(ctx context.Context, db *db.Queries) *Container {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandler(userService)

	collectionRepo := collection.NewRepository(db)
	collectionService := collection.NewService(collectionRepo)
	collectionHandler := collection.NewHandler(collectionService)

	return &Container{
		DB:                db,
		userHandler:       userHandler,
		collectionHandler: collectionHandler,
	}
}
