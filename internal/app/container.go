package app

import (
	"context"

	"github.com/Alkush-Pipania/carter-go/internal/modules/user"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

type Container struct {
	DB          *db.Queries
	userHandler *user.UserHandler
}

func NewContainer(ctx context.Context, db *db.Queries) *Container {
	userRepo := user.NewRepository(db)

	userService := user.NewService(userRepo)

	userHandler := user.NewUserHandler(userService)

	return &Container{
		DB:          db,
		userHandler: userHandler,
	}
}
